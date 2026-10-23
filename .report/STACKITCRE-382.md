# Code Review Report - STACKITCRE-382

**Rama analizada:** `feat/STACKITCRE-382-registry` (comparada contra `origin/main`)  
**Fecha:** 2026-10-23  
**Revisor:** Senior Code Reviewer  

---

## Resumen Ejecutivo

Se revisó la implementación de los nuevos comandos para Container Registry (`stackit registry artifactory [create|delete|describe|list|update]`), la configuración del cliente API y la documentación generada.

La estructura general y la integración con Cobra y el SDK de STACKIT son sólidas y siguen los patrones base de `stackit-cli`. Todos los hallazgos de código (ítems 1 al 8) han sido corregidos y validados con tests unitarios y linter. El ítem 9 (enmienda de commit) queda reservado para la etapa de commit/rebase interactivo cuando el usuario lo autorice.

---

## Estado de Resolución

- [x] **Ítem 1 (Major):** Validación conjunta obligatoria de credenciales proxy (`MarkFlagsRequiredTogether`) en `update.go`.
- [x] **Ítem 2 (Major):** Manejo seguro de secretos con `flags.SecretFlag`, `SecretFlagToStringPointer` y `json:"-"` en `create.go` y `update.go`.
- [x] **Ítem 3 (Major):** Paginación completa con seguimiento de `NextPageToken` y respeto de `--limit` en `list.go`.
- [x] **Ítem 4 (Minor):** Soporte para `--registry-custom-endpoint` en `config set`, `config unset` y `profile.json`.
- [x] **Ítem 5 (Minor):** Verificación de respuesta nil en `GetArtifactoryName` y suite de tests unitarios en `utils_test.go`.
- [x] **Ítem 6 (Minor):** Captura de entidad actualizada retornada por los `WaitHandler` en `create.go` y `update.go`.
- [x] **Ítem 7 (Suggestion):** Formato homogéneo de fecha `time.DateTime` en la tabla de `describe.go`.
- [x] **Ítem 8 (Suggestion):** Ampliación de la cobertura de pruebas unitarias en `update_test.go`.
- [ ] **Ítem 9 (Suggestion):** Corrección tipográfica en mensaje de commit `92766842` (`chroe` -> `chore`) pendiente de rebase/amend cuando el usuario lo indique.

---

## Hallazgos

### 1. Descarte silencioso de credenciales por omisión de `MarkFlagsRequiredTogether`
- **Severidad:** `Major`
- **Ubicación:** [update.go:103-108](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/update/update.go#L103-L108), [update.go:148-153](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/update/update.go#L148-L153)
- **Problema concreto:** En `update.go`, no se vinculó `cmd.MarkFlagsRequiredTogether(proxyAccessKeyFlag, proxyAccessSecretFlag)` como sí se hizo en `create.go`. Si un usuario ejecuta el comando especificando únicamente `--proxy-access-key` (o `--proxy-access-secret`), la validación pasa con éxito. No obstante, en `buildRequest` la condición `model.ProxyAccessKey != nil && model.ProxyAccessSecret != nil` se evalúa como falsa, por lo que `proxyReq.Credentials` queda `nil`. La petición se envía a la API sin credenciales y el CLI emite `"Updated artifactory..."`, descartando silenciosamente la credencial provista por el usuario sin ningún aviso o error.
- **Sugerencia de corrección:**
  Configurar la regla de flags en `configureFlags`:
  ```go
  func configureFlags(cmd *cobra.Command) {
      cmd.Flags().String(proxyUrlFlag, "", "Upstream registry URL to proxy")
      cmd.Flags().String(proxyDescriptionFlag, "", "Description of the proxy registry")
      cmd.Flags().String(proxyAccessKeyFlag, "", "Access key for upstream proxy registry authentication")
      cmd.Flags().String(proxyAccessSecretFlag, "", "Access secret for upstream proxy registry authentication")

      cmd.MarkFlagsRequiredTogether(proxyAccessKeyFlag, proxyAccessSecretFlag)
  }
  ```

---

### 2. Exposición de secretos en línea de comandos, logs de depuración y falta de soporte para `@file`/`stdin`
- **Severidad:** `Major`
- **Ubicación:** [create.go:120](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/create/create.go#L120), [update.go:107](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/update/update.go#L107), [create.go:33-41](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/create/create.go#L33-L41)
- **Problema concreto:** El flag `--proxy-access-secret` se registró como un flag string plano ordinario (`cmd.Flags().String`). En todo el CLI (`alb`, `cdn`, `load-balancer`, `intake`), los secretos deben configurarse mediante `flags.SecretFlag(proxyAccessSecretFlag, params)` y extraerse con `flags.SecretFlagToStringPointer`. La implementación actual causa:
  1. Exposición del secreto en el historial de shell y en la lista de procesos (`ps aux`).
  2. Imposibilidad de suministrar el secreto de forma segura mediante un archivo (`--proxy-access-secret @./secret.txt`) o stdin.
  3. Fuga del secreto en texto plano en los logs de depuración: `p.DebugInputModel(model)` convierte `inputModel` a JSON y luego a mapa; al ser un campo exportado sin ofuscación, el secreto se imprime en los logs de debug.
- **Sugerencia de corrección:**
  Usar `flags.SecretFlag` y proteger el modelo:
  ```go
  // En configureFlags(cmd *cobra.Command, params *types.CmdParams):
  secretFlag := flags.SecretFlag(proxyAccessSecretFlag, params)
  cmd.Flags().Var(secretFlag, proxyAccessSecretFlag, secretFlag.Usage())

  // En parseInput:
  ProxyAccessSecret: flags.SecretFlagToStringPointer(p, cmd, proxyAccessSecretFlag),

  // En inputModel:
  ProxyAccessSecret *string `json:"-"`
  ```

---

### 3. Falta de paginación en `list` (pérdida de resultados)
- **Severidad:** `Major`
- **Ubicación:** [list.go:62-73](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/list/list.go#L62-L73)
- **Problema concreto:** `ApiListArtifactoriesRequest` admite `NextPageToken` y `PageSize`, y la API responde con un `NextPageToken`. Actualmente, `list.go` solo ejecuta un único llamado:
  ```go
  req := buildRequest(ctx, model, apiClient)
  resp, err := req.Execute()
  ```
  Si el proyecto cuenta con más elementos que el tamaño de página por defecto o si el usuario indica un `--limit` superior a una sola página, los elementos posteriores nunca se obtienen ni se muestran, reportando una lista incompleta de artifactories.
- **Sugerencia de corrección:**
  Implementar el bucle de paginación siguiendo `NextPageToken` hasta alcanzar el `limit` o hasta que no haya más páginas, siguiendo el patrón de servicios como `beta/cdn/distribution/list/list.go`.

---

### 4. Inconsistencia de configuración: endpoint de registry no expuesto en `config set` / `config unset`
- **Severidad:** `Minor`
- **Ubicación:** [config.go:62](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/pkg/config/config.go#L62), [set.go:58](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/config/set/set.go#L58), [unset.go:61](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/config/unset/unset.go#L61)
- **Problema concreto:** Se añadió `RegistryCustomEndpointKey` a `internal/pkg/config/config.go` y se consume en `internal/pkg/services/registry/client/client.go`. Sin embargo, no se crearon los flags correspondientes `--registry-custom-endpoint` ni en `stackit config set` ni en `stackit config unset`, ni se agregó al template de perfiles `profile.json`. En consecuencia, ningún usuario puede configurar ni eliminar dicho endpoint mediante los comandos estándar de configuración.
- **Sugerencia de corrección:**
  Agregar el flag `registryCustomEndpointFlag = "registry-custom-endpoint"` y su enlace con viper en `internal/cmd/config/set/set.go` e `internal/cmd/config/unset/unset.go`.

---

### 5. Falta de verificación de nil en `GetArtifactoryName` y ausencia total de `utils_test.go`
- **Severidad:** `Minor`
- **Ubicación:** [utils.go:10-15](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/pkg/services/registry/utils/utils.go#L10-L15)
- **Problema concreto:** La función `GetArtifactoryName` desreferencia directamente `resp.Name` tras el llamado API sin verificar previamente si `resp != nil`. Además, el paquete `internal/pkg/services/registry/utils` no cuenta con archivo de tests unitarios (`utils_test.go`), a diferencia de prácticamente todos los demás servicios bajo `internal/pkg/services/*/utils/`.
- **Sugerencia de corrección:**
  Agregar validación de nil y crear `utils_test.go` con mocks del servicio:
  ```go
  func GetArtifactoryName(ctx context.Context, apiClient registry.DefaultAPI, projectId, region, artifactoryId string) (string, error) {
      resp, err := apiClient.GetArtifactory(ctx, projectId, region, artifactoryId).Execute()
      if err != nil {
          return "", fmt.Errorf("get artifactory: %w", err)
      }
      if resp == nil {
          return "", fmt.Errorf("get artifactory: empty response")
      }
      return resp.Name, nil
  }
  ```

---

### 6. Descarte de la respuesta actualizada del `WaitHandler` en modo sincrónico
- **Severidad:** `Minor`
- **Ubicación:** [create.go:99](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/create/create.go#L99), [update.go:88](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/update/update.go#L88)
- **Problema concreto:** Al completarse la espera (`wait.CreateArtifactoryWaitHandler` o `wait.UpdateArtifactoryWaitHandler`), el handler retorna la entidad `*registry.Artifactory` con su estado actualizado. Sin embargo, el código ignora dicho retorno (`_, err = ...`) y le pasa a `outputResult` el `resp` inicial devuelto por la llamada de creación/actualización. En formatos de salida como JSON (`-o json`), la respuesta puede reflejar estados intermedios en vez del estado consolidado tras la espera.
- **Sugerencia de corrección:**
  Capturar la respuesta del wait handler:
  ```go
  err := spinner.Run(params.Printer, "Creating artifactory", func() error {
      resp, err = wait.CreateArtifactoryWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, artifactoryId).WaitWithContext(ctx)
      return err
  })
  ```

---

### 7. Formato crudo de fecha en salida de tabla de `describe`
- **Severidad:** `Suggestion`
- **Ubicación:** [describe.go:108](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/describe/describe.go#L108)
- **Problema concreto:** `table.AddRow("CREATED AT", artifactory.CreateTime.String())` invoca el método por defecto de Go `time.Time.String()`, imprimiendo timestamps con sufijos de zona y reloj monotónico desalineados con el resto de la CLI. En toda la aplicación se utiliza `artifactory.CreateTime.Format(time.DateTime)` o `time.RFC3339`.
- **Sugerencia de corrección:**
  ```go
  table.AddRow("CREATED AT", artifactory.CreateTime.Format(time.DateTime))
  ```

---

### 8. Cobertura incompleta de casos de prueba en `update_test.go`
- **Severidad:** `Suggestion`
- **Ubicación:** [update_test.go:74-140](file:///Users/catalan/Project/STACKIT/stackit-cli/internal/cmd/registry/artifactory/update/update_test.go#L74-L140)
- **Problema concreto:** La suite de pruebas de `update` solo cubre el escenario en el que se modifica `--proxy-url`. No existen casos de prueba para verificar:
  - Actualización de `--proxy-description`.
  - Actualización con credenciales (`--proxy-access-key` y `--proxy-access-secret`).
  - Validación de combinaciones inválidas (por ejemplo, clave sin secreto).
- **Sugerencia de corrección:**
  Añadir casos de prueba para cada bandera de proxy en `TestParseInput` y `TestBuildRequest`.

---

### 9. Typo en mensaje de commit (Conventional Commits)
- **Severidad:** `Suggestion`
- **Ubicación:** Commit `92766842`
- **Problema concreto:** El commit `92766842` contiene un error tipográfico en el tipo de cambio: `chroe(cli): revert bump` en lugar de `chore(cli): revert bump`.
- **Sugerencia de corrección:**
  Corregir mediante rebase interactivo (`git commit --amend`) antes de integrar a la rama principal.

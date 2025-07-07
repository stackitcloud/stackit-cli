package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

var (
	testProjectId = uuid.NewString()
)

const (
	testClusterName    = "test-cluster"
	existingKubeConfig = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURCVENDQWUyZ0F3SUJBZ0lJSjFTZ1NWTjhnMmt3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TlRBeE1UQXhNakkxTlRSYUZ3MHpOVEF4TURneE1qTXdOVFJhTUJVeApFekFSQmdOVkJBTVRDbXQxWW1WeWJtVjBaWE13Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLCkFvSUJBUUM4ZXIwam1aS05STlR6Z2dCV3Q1cXMvaW94NXkxY2xBMHBGRHYwOWNmMGtmVGRVQWE3bmpqU0F2WlYKVFpsQlFFaW40Um9PTm1TZzdVMzVWN3FMSW56UVNmZXFuYi9wK05pODhDbkZvMThleUVnb3pHQklTTFpHK0EybQpuNFFEV3k3bVV1UUxFRnpjNjFpazdBQ0F5akZwRDlVdkdSdkxxVGJTQWcwYitYbktqbUUyWVgzTnRLbnJWOUN0CktrTG83K2JSa0MyemNkVnlraExhODhaR1BORUhjdVp2Uk0zQW5NclVGdGVvc0Fjb09xVW4xK09mYlhwUUlsTC8KKzBvRjcwN09Vc2tOUit0WEp4Z1VXL1R4Q0lONTYwU2E4eDVlWjB2VTZNR3ZOSTYwZ3h2S1lGL0pKa0pxU0NwNQovWWhpVmZ2QnNOSG5tVUZsNEdpOGFVMFNVTjRiQWdNQkFBR2pXVEJYTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQCkJnTlZIUk1CQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJTUlkxVVhOamlMbFJLWktuSHJWRU55djA4aUp6QVYKQmdOVkhSRUVEakFNZ2dwcmRXSmxjbTVsZEdWek1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQVFaOVcwVFdvMQp4UFhPZU9xWHV6aFgzSkRoY0JVRkZyUVlOcHBMSmtqOWdUVm5Eck16b1dmeW9FQXRtT1ZQWURuTnEyTFhOSnpmClltd3RiUGxPemhGYkpWZVBWR0tLZktrUXZ1K3BhZGRtUHRhTzdUcnZqblRHeDhXczJadE5xK20wbkRGRUN4SDkKc1o2K1IycWhBUWNnSGdQWFZQdTdxSXFmbkNWRDkyeGprTE40c2JLZjRMb2x0R3hZbTBTWVZuY09rTFlBL3BvawpqTCsvODRJQXRrRXlEL21VdVF4MEsyVzFvVUM4dDRyMUlPZ3Y3OHZQMkRDRlBuZDVvbTJBM1dCNHY2dUFNZWc0Cnk3Y3FTcjBlSzJhNFQvMUtpTEdzYXI1V01ONTNwMjFiOGJMSTlISGNJMkh6c0tOdEdpNGFOT0hsWWkwUFgrUW0KT3U4NW4ycVdwSUxmCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://127.0.0.1:61274
  name: existing-cluster
contexts:
- context:
    cluster: existing-cluster
    user: existing-cluster
  name: existing-cluster
current-context: existing-cluster
kind: Config
preferences: {}
users:
- name: existing-cluster
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURLVENDQWhHZ0F3SUJBZ0lJYWFEL3lTemlKM1F3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TlRBeE1UQXhNakkxTlRSYUZ3MHlOakF4TVRBeE1qTXdOVFJhTUR3eApIekFkQmdOVkJBb1RGbXQxWW1WaFpHMDZZMngxYzNSbGNpMWhaRzFwYm5NeEdUQVhCZ05WQkFNVEVHdDFZbVZ5CmJtVjBaWE10WVdSdGFXNHdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEQSs2WjkKWU02RC9DK2VNWnJQRHZoR0VIRk4zeDVXdFMrVWlsb1F3QkJBSXdUNXFQczVnSERWK1cyWjdjT3VGNVFEYlpyUQo3dktWSUtlWXQ0Mk9SZytYQktibHhDV1VpdFZDdmZZbHJYKzlaY0JGL2dFaVBjOE9aK2h0Q1pPNlgyZ3d0WVNOCkgwZ1lLOTlhOFRWUWxlWm9Eem93WlE0Um5aSjhkRGo1STA2blRjdkk3bDBlMWt3VnM5aXFLRHpyekRhYnhqb0EKamZkcUpiZTVkOFc0ZTloTTRBdVRUbFRkWmFVTWFnUHhyaWxEOU9mUXhaUmlReFIzNkhSOHZabm9TcndXeWh5ZApqall0TFQvcE00UXAybUU5NFJqVWE2ekNUVlJKeWduY3RHVnpDRi84RDc1TVU4OVhmVjltQVV5L3BoR1M5MDdjCjlXbzE4Um42TytHNHYwdFRBZ01CQUFHalZqQlVNQTRHQTFVZER3RUIvd1FFQXdJRm9EQVRCZ05WSFNVRUREQUsKQmdnckJnRUZCUWNEQWpBTUJnTlZIUk1CQWY4RUFqQUFNQjhHQTFVZEl3UVlNQmFBRkpGalZSYzJPSXVWRXBrcQpjZXRVUTNLL1R5SW5NQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUMySVRwUlM1SlU3bGpkeDVRMlkwQzBkZG8yCm9PSmp5TUhVQXJ5ZTIyM2xOd1R1OTNHZXkvUjNIOHNpYWxDRURXdFR0cCsrY3BucW1ON05ia3UvWFI5SUlFdlIKYTNZS3VvbGdOTGtLaEtqMWQ0NVAxeEs0VE5CV1hSV2FMbksxcTdLVWxWWHp2bjdSN3RDY0NtNk90S3d4OUl2WgorRGhUU0pobFEzTVNmNXhjMUdOMm9qb0pPWmVlOXFNc3R1RzdPUVl1M08yUitYVUIwRHgzNnlPeFR2S0NBZ24xCm55Yk5FS0Nia1BmTXdvSU5aTm9iSWE3Y2VHcTdOMzRHaCs3Vi9iazUrQmhoTzVJRTRPeDYvUUxQc1B2ZGtOZHcKSkFyclQ3QytHSkF1UzNXQ2dYUXRyRWFyT3drWHhqajFPc3NuNjdMNlpONG01SkYzWHViSmdQUGZ3L2NECi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBd1B1bWZXRE9nL3d2bmpHYXp3NzRSaEJ4VGQ4ZVZyVXZsSXBhRU1BUVFDTUUrYWo3Ck9ZQncxZmx0bWUzRHJoZVVBMjJhME83eWxTQ25tTGVOamtZUGx3U201Y1FsbElyVlFyMzJKYTEvdldYQVJmNEIKSWozUERtZm9iUW1UdWw5b01MV0VqUjlJR0N2Zld2RTFVSlhtYUE4Nk1HVU9FWjJTZkhRNCtTTk9wMDNMeU81ZApIdFpNRmJQWXFpZzg2OHcybThZNkFJMzNhaVczdVhmRnVIdllUT0FMazA1VTNXV2xER29EOGE0cFEvVG4wTVdVCllrTVVkK2gwZkwyWjZFcThGc29jblk0MkxTMC82VE9FS2RwaFBlRVkxR3Vzd2sxVVNjb0ozTFJsY3doZi9BKysKVEZQUFYzMWZaZ0ZNdjZZUmt2ZE8zUFZxTmZFWitqdmh1TDlMVXdJREFRQUJBb0lCQVFDS0lWWFk5anE3VS8zTgpjRm9MalA1K1AvU3B0V01rMHdsY2UrN2RnR3ZoVEcrYU42NmlTT0g2OWs3UjE5S3hRS1VzRXY2MlArSVloY2dRClVvbWE1V0R4U2w0ZnBkYjBUSzg2MTNkaEhwK0pORlI4aE1QUSs0YkNHL1BNWUFlQ1poblFpNHgxNm9jUzdnd3cKTHVoblp2UUZWYWpqek9GV0VJQXlYb29OSVkyQng3bjlzRlBGYmZSK1NOVVhuWHNHemFkMlArVmIyTkFCUjRFLwp4K2dYWlhFKzFnU0RhK25ZVHBiaG1hd3hreStEQnZBQlRWTzlWY2J2ZWoybDZ2WjAwK2lMTm9rYjF1UmJmbzNECkdEN2RZTjRYdCtwWXRMdFJYRGNqb2Q3OXpFcmJ4UkE4ZWoxblllOFpXQUNZa0ZOT3lpRHlJY3dFbWtDNXhlcHAKS1ByRGVCeEpBb0dCQU5XYzI4cFY4SDhRWm1Hb25QQkNZUUNrY2NLYnpEaXpwa0ZKMlZNVXZ5TG1Ia0w5bWlWUApQb1RsdXF4T2htMHhyRlNRaEFTQUlUaG0rWHN0c0pYdjNSd2dIZVdadTluUEVPeWpRcG02bTNEa1ZVK29kdTRGCnYwa25qdlduUTRPZnVQeDlCV3UyN1I3d1VBNHBqNUk0MGtlMVovdDZwdzZjeWFBckZ5L01HODZmQW9HQkFPZEcKMXRocFNUT3dZbEltWWoxNDdTSVJyb0VaSTNSaUNBVUh1ME54VEJObk5WL1JNVDdaNGVpZkRMMndXc0s1Q1Y0aQpFR2hBODRxYVB0dTFCaVhwTmdpMDBBdllWUGN6d0VDa3hocFdBeTJVRGZSc2FENnNYQ04ycVdtcGdjQzBTOWpICkdqUkdnSVFselhWNFVVcHFTYzVEUDZBYUFzRkhxVU1aT3dRZTgyck5Bb0dCQU5FN2FLbml6Y09ZQzhDQ2lONXAKRmx5cnRtWVpkc3JmWk95MGFqT1BzYng4VEkzdm04b0p1Y0l3eDAwNVNVQ3hsQXZzMWZNV2tmT09JYlkreGFYSApvZnVIbGVFc1dTejZQcWliTFlRb25WTFJ4S0pXNzg4clAvZG0wUWZiZ3l6dENTUC9UWXo1UzMrdmdhcXRtTnh2CjNjQ3hkcDJEd1JoMkNLUmpNTDMzbmhFZkFvR0FDNmNRRUJ0TjZ1TEtNV1Zwc2JzMEIzRm9uMnlLMHNSVnJ4c3kKbmpWSkpma2ZRVktpN28yL3loNnBYNjFSQlZxWlZEclhKTW1RKzd6RnlnQVc3VFlRMk9OelVBVjRVblF6RFk2Lwp4SGZzOVJEdW14QVRPSVVxcDBiRlJtT1ovQUdaaUxTUFoyN2Q3c3FRelloZ1lDVjJ6b09vNHdJc2ZWeUU5TEtDCnZMUnFnMGtDZ1lFQXlJRUdjeHQxcTIwdUhYUTFLTU92V2xWUUJCQklPUUJjeXoyR0djcWFGOHhSKzJCOGc3R2YKbEh4dHBvaTNNQUxTVXlhOTQzZEpMUHA4Q0xSOTBkQWtqZ1JROURPN2wyYWlWYWVncTA0NURCMnBwN05YVlc4NgptUXFPZUJRYzcyY0ZYdk9YZmRKUUQwME5HZThlS0VjTWN2QlhxTVIrSUtEdGozcGlKVjlsSHpBPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`
	newKubeConfig = `apiVersion: v1
clusters:
  - cluster:
      certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURCVENDQWUyZ0F3SUJBZ0lJTjAvdmZkM3RCeGd3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TlRBeE1URXhOekV4TVRkYUZ3MHpOVEF4TURreE56RTJNVGRhTUJVeApFekFSQmdOVkJBTVRDbXQxWW1WeWJtVjBaWE13Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLCkFvSUJBUURJS25lRWZrM0F0WWlhanZyYWcwdU1zZUd1Q3BuZW80OXl3T0NFSmF0ZnVncVZXVXJ3cVd1WEdjUVgKTWp5MTZEVGxlR2YxeS83NXJuRUY1cld0Vm5wMDlNc0w1NW5YM0ZnT21SY3ozNmxtYTBOMmdMQU5RR0VmZU50NQpsa0Y5R2t6VFZMVy84alNWcXRkaTBCTm8xejEya0FCUm5yM1M0bWU0cExma0xFeWZKQTFQcnlpVUp0NnFBbldrCkUwV2RxbmJJMGRHQWZpZ3hTVFRZK09PMExWbjdJaG1QTGpPVEhHb0JRaW1DL091ZEZFK01FZG1kQkNOTHgzeE4KRDlSbk1taUxjVkVlSDlvVTFjYUdRamRIbXhnRUpJbStTOVdmWDZuRSsxOUpDZ0dkTS9KaFVtT0xRQWg4NzhMcQptc085WlNYdXFweW9ROTBhRDBDaFNNdzJyOXBQQWdNQkFBR2pXVEJYTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQCkJnTlZIUk1CQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJRMXRjTE5rMmVjRkFJRDl5citZMnUyaHI4OWJEQVYKQmdOVkhSRUVEakFNZ2dwcmRXSmxjbTVsZEdWek1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQWdUZGJkTzZQNQo2M2hiZVRsS1E2UkpzRlkrdUdIeXcyMXNGU205Ni9vblZhOS91SjNQZ3BsMndKaFhpanZmZnNQamg2ekpkdTJXCll4WWkxcHdEWGZtMHpsNHJQMEcwQmkzL2Y1VkU0dkRnSmUwcDRKdkx2MWVmclZBcGhpakJiRkFHVTh6WVVPdEUKM2pGNy92ZDkvVUwxRWwzNVNRZjdEWWJhQ2NndzByS0tiNkQwaUZJcjJCRFZqbE01VDhqRzdETEk0a3pXTzFaTQpmNHh4ay9MQjBpY1R0a1RVRGQzcjBtZmFzNUdqR0lDR2QzbUpHbWY3bzFScXVyYlZ3dmVPWE5oL2tud2hnNGZqCitsTjJvaHpuaWdkTVNNQ1FnbDQ2NlowQTZvVDUrNUV6a2JwYS8yRDQ1cVN0ZGZBbTNtQ0RhdHdUelc1RlBudFMKMm0weVo2ZWVydkE4Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
      server: https://127.0.0.1:55209
    name: my-new-super-ske-cluster
contexts:
  - context:
      cluster: my-new-super-ske-cluster
      user: my-new-super-ske-cluster
    name: my-new-super-ske-cluster
current-context: my-new-super-ske-cluster
kind: Config
preferences: {}
users:
  - name: my-new-super-ske-cluster
    user:
      client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURLVENDQWhHZ0F3SUJBZ0lJUmpoS0w0dlJWSFV3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TlRBeE1URXhOekV4TVRkYUZ3MHlOakF4TVRFeE56RTJNVGRhTUR3eApIekFkQmdOVkJBb1RGbXQxWW1WaFpHMDZZMngxYzNSbGNpMWhaRzFwYm5NeEdUQVhCZ05WQkFNVEVHdDFZbVZ5CmJtVjBaWE10WVdSdGFXNHdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEVE5WSmEKekJHZWU4OXVRNjVZWEdhT1pwTWJTZE9tcWFyNUlVbkRTUEpMbHdJKzkyWVRrcFBKcXFncWEwa2FZYVdZUmFlTQpCNVlDeTRpNjNXSTBYYlgvMW9LNUFPZ2xXL1FwcGczWnc5K3ZPYXdtdEpqUHQ1T2xEVWRONGdmYm40TjV1OWpoCmltQ09wak5VL285NzNZZy9nM3pqNi9nUm9EYldhaW5wSDltTk1nOHFTS0xaNkNpUlp2VjZuYkgyVDVSa3ZVVWgKUDNWN09CZE1oUlp3MW1rVVRQVXY5T056VVBubFFaS3hwWXphYjBiZm92eFd6UDhxQkVIdk9xaXZoWFhaaGp1bApaTU1OMjYrN2RyS3lCWS8rRnBmeGpqb3AyZytUSlMxNHhhOTh0dCtqT3dUUkI5aWh1WUQzTnlVbEZXVjhiUG51CnJqSW52ckxVcjkvQzB2cmhBZ01CQUFHalZqQlVNQTRHQTFVZER3RUIvd1FFQXdJRm9EQVRCZ05WSFNVRUREQUsKQmdnckJnRUZCUWNEQWpBTUJnTlZIUk1CQWY4RUFqQUFNQjhHQTFVZEl3UVlNQmFBRkRXMXdzMlRaNXdVQWdQMwpLdjVqYTdhR3Z6MXNNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUFZQkpld0ZwMTJnbkxQM1hGQ09JaXRZZWVnCkVmMjQwLysvaVFUUXQreHNjTU1ITGF4VjNFNEgxZ3JyNDdXUjE0bDdlbE1ING5qWnZzU3djSUZsa1RieVR6eW0KeW9XamhQQ0M2WWpzZHFEM2Vlc1ZpV2xhZkthczFrNmtmWHhVR2EvSUtQNzJoQ2tub2pia2o3amlSdjgrMTd5NgpKa2JIaXNYLzFqM2R1VHVIdDNORXJnNmNud0M5MGlldjZFZVFaV0oxaG5NSHhDMkRYMEdvOW14ZDlPYWFVODdBCkhBNDMzRnVJQWpoZjRWN2Vma3dGQU1ZMEhZSjZQaFZqTXdNWmdKczhLSHhVdjl3Y0xYMlFPUC9TSmhRZUtMV1UKYTFHTWlzTFBNc2NmL2JjU051SVpxMTR5S0xSelEwL1FIUW1PVVdSZDIva002MmxhbFl5Rlk2V0J4cCt3Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
      client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBMHpWU1dzd1JubnZQYmtPdVdGeG1qbWFURzBuVHBxbXErU0ZKdzBqeVM1Y0NQdmRtCkU1S1R5YXFvS210SkdtR2xtRVduakFlV0FzdUl1dDFpTkYyMS85YUN1UURvSlZ2MEthWU4yY1Bmcnptc0pyU1kKejdlVHBRMUhUZUlIMjUrRGVidlk0WXBnanFZelZQNlBlOTJJUDROODQrdjRFYUEyMW1vcDZSL1pqVElQS2tpaQoyZWdva1diMWVwMng5aytVWkwxRklUOTFlemdYVElVV2NOWnBGRXoxTC9UamMxRDU1VUdTc2FXTTJtOUczNkw4ClZzei9LZ1JCN3pxb3I0VjEyWVk3cFdURERkdXZ1M2F5c2dXUC9oYVg4WTQ2S2RvUGt5VXRlTVd2ZkxiZm96c0UKMFFmWW9ibUE5emNsSlJWbGZHejU3cTR5Sjc2eTFLL2Z3dEw2NFFJREFRQUJBb0lCQVFDRDBQV1RJV1dscWRQdQpGMk9LVmpEVGt3VWd0TlRaWVc4SmlWTUdCRkxrQmwwcWV6RkQ2ZWsrcGJuS3I2YXlSbHNaUysram4yQnFZaWoxCnB4R1JhU01iaHYrVEF4UGZyU0lYbEVGMHRhQzNOYUZSanNrSWFxUkZFS0o5NHlIUVdoK3VMQ1RScnBGUXRqMjMKUUNEQXg2UXZMNXNVak1NSURSdnNlZG1xVzJ4bGg4UkF5RUdYVi9sUmJ5ZTdEOTIrWVpwd21kV3dsa2tiZy8yTQowdHF1R1k0Qk1XTFY0K09DVlNmVWVEWU1nZkZIL0RVWThUdUIvNitzVm9rUnhLalhYbjYzN1c4Q2dJWUVaQngrCkE5TG8vYk1YN0RaSDRmS0RyRCsycVQ1SDNUTDFIc3BtSXJ1Mi84RllCZ08ySjNzZVdHdHdtelVXalVzL2ExSGoKdXZMamNCTjVBb0dCQU83YitESTBsdFRGT29MSERISGNZdXZqMTYydU96bk51ejNXa1R0Sng3QzZJSVpVd2YwSQpuM2pJWXhKRi9yVVZUZzZPbU5XNXpGdDA2QTVWQitwZ2RNblFhOHMybVNldFlKVW51eE42emRsOGJoblZ6dXUzCi8walM3cU1pWGg5aU8vRlZ6VDNxcnNqU1VnMmNCRTA3WlZweU0ycVNMUlkrVmdiRDg4aUdUbXhiQW9HQkFPSmQKWWVNc1JpVVZ5Wk5sZU1ra3puS2pjYXoyOE9Vb0NyZjd0dVhaYUpqRDdWZncyWmNBd0cvZG5lZ3M2YmEvck54bgplMXU3Rm05VlNTR2pNejJEaC9QdlNuQlZReGtQeHo1ZFRja2V0RUJSQk1XaVV1enI2UUFXdmZudEZXcWNZTkpvClBCVWY3c2k4Wk1rMjJpanR1OWxEVnRRUFpJdDZUMzJrb0Z3eHNrcHpBb0dBYjQ0c2pNWWk2NXh4aDBLUGZWNEEKbFVzRUlBbVBmNSttSTJ0aXlOM2NkWjE0TTBUQ2xQckNBQmNXcmlJaW8xQWY5SXlFdE16aHRKVVZEQnlLWmR4RwpyenE0SFdDU2h3Vmlaa2I0Q0ZFQ2N1QzZTemFnUFZiaDA1RXdBdUM2Tk00Y1VNcFI0T2tLV0tCaDBobGJxUFprCmo2bG1lZzlySDBoZHhTc2ZZRGZaeUtFQ2dZQnZZMVk4ekZlRC9qR2YxMG5WYU1neC94MTc2RlBuMzRsT3VZMXAKazA3MkJVdHdmN01DckRzRmtQOFg5YW5YNUgveVFQV2gwUEVjUGRKcnUvd0Y1QWh0VDYzSWt4d2VZL1krU1BseQo0eW45a0NDU0ErdGNiRVhPWm1KN2JsK2dnMnpkZks4OEVlZVZYYWNXb0dnL3hhUXZLQVM4K3dvVjNFenJYYXdQClVlRVM0d0tCZ1FEUm9QbXkvNloySUdERkRReWt3YmFMRDlvQlZqN3BJSTI0NmlLM1hwQmRtRGFVR0hLYnRiNmUKYXNYRWNQQmp0enYvTzVOM2dlZWFYREduaW5XcXJJZm1FTzIyMDhmQ0VCc0RWc3RQMDhxRnorekFSMnJEQm9xbQpFVkwxN0o0Q2J6Tlh4bStOT1R6aVhCN2tLVWhNQUFBbmkwcXQ1QXN0QlJpcENuMER4Y2JpekE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`
)

type skeClientMocked struct {
	listClustersFails        bool
	listClustersResp         *ske.ListClustersResponse
	listProviderOptionsFails bool
	listProviderOptionsResp  *ske.ProviderOptions
}

const testRegion = "eu01"

func (m *skeClientMocked) ListClustersExecute(_ context.Context, _, _ string) (*ske.ListClustersResponse, error) {
	if m.listClustersFails {
		return nil, fmt.Errorf("could not list clusters")
	}
	return m.listClustersResp, nil
}

func (m *skeClientMocked) ListProviderOptionsExecute(_ context.Context, _ string) (*ske.ProviderOptions, error) {
	if m.listProviderOptionsFails {
		return nil, fmt.Errorf("could not list provider options")
	}
	return m.listProviderOptionsResp, nil
}

func TestClusterExists(t *testing.T) {
	tests := []struct {
		description      string
		getClustersFails bool
		getClustersResp  *ske.ListClustersResponse
		isValid          bool
		expectedExists   bool
	}{
		{
			description:     "cluster exists",
			getClustersResp: &ske.ListClustersResponse{Items: &[]ske.Cluster{{Name: utils.Ptr(testClusterName)}}},
			isValid:         true,
			expectedExists:  true,
		},
		{
			description:     "cluster exists 2",
			getClustersResp: &ske.ListClustersResponse{Items: &[]ske.Cluster{{Name: utils.Ptr("some-cluster")}, {Name: utils.Ptr("some-other-cluster")}, {Name: utils.Ptr(testClusterName)}}},
			isValid:         true,
			expectedExists:  true,
		},
		{
			description:     "cluster does not exist",
			getClustersResp: &ske.ListClustersResponse{Items: &[]ske.Cluster{{Name: utils.Ptr("some-cluster")}, {Name: utils.Ptr("some-other-cluster")}}},
			isValid:         true,
			expectedExists:  false,
		},
		{
			description:      "get clusters fails",
			getClustersFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &skeClientMocked{
				listClustersFails: tt.getClustersFails,
				listClustersResp:  tt.getClustersResp,
			}

			exists, err := ClusterExists(context.Background(), client, testProjectId, testRegion, testClusterName)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if exists != tt.expectedExists {
				t.Errorf("expected exists to be %t, got %t", tt.expectedExists, exists)
			}
		})
	}
}

func fixtureProviderOptions(mods ...func(*ske.ProviderOptions)) *ske.ProviderOptions {
	providerOptions := &ske.ProviderOptions{
		KubernetesVersions: &[]ske.KubernetesVersion{
			{
				State:   utils.Ptr("supported"),
				Version: utils.Ptr("1.2.3"),
			},
			{
				State:   utils.Ptr("supported"),
				Version: utils.Ptr("3.2.1"),
			},
			{
				State:   utils.Ptr("not-supported"),
				Version: utils.Ptr("4.4.4"),
			},
		},
		MachineImages: &[]ske.MachineImage{
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("1.2.3"),
						Cri: &[]ske.CRI{
							{
								Name: ske.CRINAME_DOCKER.Ptr(),
							},
							{
								Name: ske.CRINAME_CONTAINERD.Ptr(),
							},
						},
					},
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("3.2.1"),
						Cri: &[]ske.CRI{
							{
								Name: ske.CRINAME_DOCKER.Ptr(),
							},
							{
								Name: ske.CRINAME_CONTAINERD.Ptr(),
							},
						},
					},
				},
			},
			{
				Name: utils.Ptr("not-flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("4.4.4"),
						Cri: &[]ske.CRI{
							{
								Name: ske.CRINAME_CONTAINERD.Ptr(),
							},
						},
					},
				},
			},
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("4.4.4"),
					},
				},
			},
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("not-supported"),
						Version: utils.Ptr("4.4.4"),
						Cri: &[]ske.CRI{
							{
								Name: ske.CRINAME_CONTAINERD.Ptr(),
							},
						},
					},
				},
			},
			{
				Name: utils.Ptr("flatcar"),
				Versions: &[]ske.MachineImageVersion{
					{
						State:   utils.Ptr("supported"),
						Version: utils.Ptr("4.4.4"),
						Cri: &[]ske.CRI{
							{
								Name: ske.CRINAME_DOCKER.Ptr(),
							},
						},
					},
				},
			},
		},
	}
	for _, mod := range mods {
		mod(providerOptions)
	}
	return providerOptions
}

func fixtureGetDefaultPayload(mods ...func(*ske.CreateOrUpdateClusterPayload)) *ske.CreateOrUpdateClusterPayload {
	payload := &ske.CreateOrUpdateClusterPayload{
		Extensions: &ske.Extension{
			Acl: &ske.ACL{
				AllowedCidrs: &[]string{},
				Enabled:      utils.Ptr(false),
			},
		},
		Kubernetes: &ske.Kubernetes{
			Version: utils.Ptr("3.2.1"),
		},
		Nodepools: &[]ske.Nodepool{
			{
				AvailabilityZones: &[]string{
					"eu01-3",
				},
				Cri: &ske.CRI{
					Name: ske.CRINAME_CONTAINERD.Ptr(),
				},
				Machine: &ske.Machine{
					Type: utils.Ptr("b1.2"),
					Image: &ske.Image{
						Version: utils.Ptr("3.2.1"),
						Name:    utils.Ptr("flatcar"),
					},
				},
				MaxSurge:       utils.Ptr(int64(1)),
				MaxUnavailable: utils.Ptr(int64(0)),
				Maximum:        utils.Ptr(int64(2)),
				Minimum:        utils.Ptr(int64(1)),
				Name:           utils.Ptr("pool-default"),
				Volume: &ske.Volume{
					Type: utils.Ptr("storage_premium_perf2"),
					Size: utils.Ptr(int64(50)),
				},
			},
		},
	}
	for _, mod := range mods {
		mod(payload)
	}
	return payload
}

func TestGetDefaultPayload(t *testing.T) {
	tests := []struct {
		description              string
		listProviderOptionsFails bool
		listProviderOptionsResp  *ske.ProviderOptions
		isValid                  bool
		expectedOutput           *ske.CreateOrUpdateClusterPayload
	}{
		{
			description:             "base",
			listProviderOptionsResp: fixtureProviderOptions(),
			isValid:                 true,
			expectedOutput:          fixtureGetDefaultPayload(),
		},
		{
			description:              "get provider options fails",
			listProviderOptionsFails: true,
			isValid:                  false,
		},
		{
			description: "no Kubernetes versions 1",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.KubernetesVersions = nil
			}),
			isValid: false,
		},
		{
			description: "no Kubernetes versions 2",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.KubernetesVersions = &[]ske.KubernetesVersion{}
			}),
			isValid: false,
		},
		{
			description: "no supported Kubernetes versions",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.KubernetesVersions = &[]ske.KubernetesVersion{
					{
						State:   utils.Ptr("not-supported"),
						Version: utils.Ptr("1.2.3"),
					},
				}
			}),
			isValid: false,
		},
		{
			description: "no machine images 1",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{}
			}),
			isValid: false,
		},
		{
			description: "no machine images 2",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = nil
			}),
			isValid: false,
		},
		{
			description: "no machine image versions 1",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{
					{
						Name:     utils.Ptr("image-1"),
						Versions: nil,
					},
				}
			}),
			isValid: false,
		},
		{
			description: "no machine image versions 2",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{
					{
						Name:     utils.Ptr("image-1"),
						Versions: &[]ske.MachineImageVersion{},
					},
				}
			}),
			isValid: false,
		},
		{
			description: "no supported machine image versions",
			listProviderOptionsResp: fixtureProviderOptions(func(po *ske.ProviderOptions) {
				po.MachineImages = &[]ske.MachineImage{
					{
						Name: utils.Ptr("image-1"),
						Versions: &[]ske.MachineImageVersion{
							{
								State:   utils.Ptr("not-supported"),
								Version: utils.Ptr("1.2.3"),
							},
						},
					},
				}
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &skeClientMocked{
				listProviderOptionsFails: tt.listProviderOptionsFails,
				listProviderOptionsResp:  tt.listProviderOptionsResp,
			}

			output, err := GetDefaultPayload(context.Background(), client, testRegion)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			diff := cmp.Diff(output, tt.expectedOutput)
			if diff != "" {
				t.Fatalf("Output is not as expected: %s", diff)
			}
		})
	}
}

func TestConvertToSeconds(t *testing.T) {
	tests := []struct {
		description    string
		expirationTime string
		isValid        bool
		expectedOutput string
	}{
		{
			description:    "seconds",
			expirationTime: "30s",
			isValid:        true,
			expectedOutput: "30",
		},
		{
			description:    "minutes",
			expirationTime: "30m",
			isValid:        true,
			expectedOutput: "1800",
		},
		{
			description:    "hours",
			expirationTime: "30h",
			isValid:        true,
			expectedOutput: "108000",
		},
		{
			description:    "days",
			expirationTime: "30d",
			isValid:        true,
			expectedOutput: "2592000",
		},
		{
			description:    "months",
			expirationTime: "30M",
			isValid:        true,
			expectedOutput: "77760000",
		},
		{
			description:    "leading zero",
			expirationTime: "0030M",
			isValid:        true,
			expectedOutput: "77760000",
		},
		{
			description:    "invalid unit",
			expirationTime: "30x",
			isValid:        false,
		},
		{
			description:    "invalid unit 2",
			expirationTime: "3000abcdef",
			isValid:        false,
		},
		{
			description:    "invalid unit 3",
			expirationTime: "3000abcdef000",
			isValid:        false,
		},
		{
			description:    "invalid time",
			expirationTime: "x",
			isValid:        false,
		},
		{
			description:    "empty",
			expirationTime: "",
			isValid:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := ConvertToSeconds(tt.expirationTime)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if *output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, *output)
			}
		})
	}
}

func TestMergeKubeConfig(t *testing.T) {
	tests := []struct {
		description        string
		location           string
		kubeconfig         string
		existingKubeconfig string
		isValid            bool
		isLocationDir      bool
		isLocationEmpty    bool
		expectedErr        string
	}{
		{
			description: "base",
			location:    filepath.Join("base", "config"),
			kubeconfig:  newKubeConfig,
			isValid:     true,
		},
		{
			description:     "empty location",
			location:        "",
			kubeconfig:      newKubeConfig,
			isValid:         false,
			isLocationEmpty: true,
		},
		{
			description:   "path is only dir",
			location:      "only_dir",
			kubeconfig:    newKubeConfig,
			isValid:       false,
			isLocationDir: true,
		},
		{
			description: "empty kubeconfig",
			location:    filepath.Join("empty", "config"),
			kubeconfig:  "",
			isValid:     false,
		},
		{
			description:        "kubeconfig bad content",
			location:           filepath.Join("empty", "config"),
			existingKubeconfig: "hola",
			kubeconfig:         "kubeconfig",
			isValid:            false,
		},
		{
			description:        "kubeconfig content",
			location:           filepath.Join("content", "config"),
			kubeconfig:         newKubeConfig,
			existingKubeconfig: existingKubeConfig,
			isValid:            true,
		},
	}

	baseTestDir := "test_data/"
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testLocation := filepath.Join(baseTestDir, tt.location)
			// make sure empty case still works
			if tt.isLocationEmpty {
				testLocation = ""
			} else if tt.existingKubeconfig != "" {
				dir := filepath.Dir(testLocation)

				err := os.MkdirAll(dir, 0o700)
				if err != nil {
					t.Errorf("error create config directory: %s (%s)", dir, err.Error())
				}

				err = os.WriteFile(testLocation, []byte(tt.existingKubeconfig), 0o600)
				if err != nil {
					t.Errorf("could not write file: %s", tt.location)
				}
				defer func() {
					err := os.Remove(testLocation)
					if err != nil {
						t.Errorf("could not deleete file: %s", tt.location)
					}
				}()
			}
			// filepath Join cleans trailing separators
			if tt.isLocationDir {
				testLocation += string(filepath.Separator)
			}

			err := MergeKubeConfig(testLocation, tt.kubeconfig)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input %s", err)
			}

			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}

			if tt.isValid {
				kubeConfigFinal, err := clientcmd.LoadFromFile(testLocation)
				if err != nil {
					t.Errorf("error loading final kubeconfig: %s", err)
				}

				kubeConfigNew, err := clientcmd.Load([]byte(tt.kubeconfig))
				if err != nil {
					t.Errorf("error loading new kubeconfig: %s", err)
				}

				// check new kubeconfig is still there
				for name := range kubeConfigNew.AuthInfos {
					_, exits := kubeConfigFinal.AuthInfos[name]
					if !exits {
						t.Errorf("the user %s does not exist in the final kubeconfig", name)
					}
				}
				for name := range kubeConfigNew.Contexts {
					_, exits := kubeConfigFinal.Contexts[name]
					if !exits {
						t.Errorf("the context %s does not exist in the final kubeconfig", name)
					}
				}
				for name := range kubeConfigNew.Clusters {
					_, exits := kubeConfigFinal.Clusters[name]
					if !exits {
						t.Errorf("the cluster %s does not exist in the final kubeconfig", name)
					}
				}

				if tt.existingKubeconfig != "" {
					kubeConfigExisting, err := clientcmd.Load([]byte(tt.existingKubeconfig))
					if err != nil {
						t.Errorf("error loading existing kubeconfig: %s", err)
					}

					// check exiting kubeconfig is still there
					for name := range kubeConfigExisting.AuthInfos {
						_, exits := kubeConfigFinal.AuthInfos[name]
						if !exits {
							t.Errorf("the user %s does not exist in the final kubeconfig", name)
						}
					}
					for name := range kubeConfigExisting.Contexts {
						_, exits := kubeConfigFinal.Contexts[name]
						if !exits {
							t.Errorf("the context %s does not exist in the final kubeconfig", name)
						}
					}
					for name := range kubeConfigExisting.Clusters {
						_, exits := kubeConfigFinal.Clusters[name]
						if !exits {
							t.Errorf("the cluster %s does not exist in the final kubeconfig", name)
						}
					}
				}
			}
		})
	}
	// Cleanup
	err := os.RemoveAll(baseTestDir)
	if err != nil {
		t.Errorf("failed cleaning test data")
	}
}

func TestGetDefaultKubeconfigPath(t *testing.T) {
	tests := []struct {
		description string
	}{
		{
			description: "base",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			output, err := GetDefaultKubeconfigPath()

			if err != nil {
				t.Errorf("failed on valid input")
			}
			userHome, err := os.UserHomeDir()
			if err != nil {
				t.Errorf("could not get user home directory")
			}
			if output != filepath.Join(userHome, ".kube", "config") {
				t.Errorf("expected output to be %s, got %s", filepath.Join(userHome, ".kube", "config"), output)
			}
		})
	}
}

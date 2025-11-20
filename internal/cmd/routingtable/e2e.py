import subprocess
import sys
import yaml
from datetime import datetime

# Static variables
PROJECT_ID = "xxx"
NETWORK_AREA_ID = "xxx"
ORG_ID = "xxx"

# Dynamic variables initialized during test flow
NETWORK_ID = ""
ROUTING_TABLE_ID = ""
ROUTING_TABLE_ID_2 = ""
ROUTE_ID = ""

def log(msg: str):
    print(f"[{datetime.now().strftime("%Y-%m-%d %H:%M:%S")}] {msg}", file=sys.stdout)

def run_command(description: str, _expected: str, *args):
    log(f"{description}")
    result = subprocess.run(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)

    if result.returncode == 0:
        log(f"Command succeeded: {description}")
        if result.stdout.strip():
            print("STDOUT:")
            print(result.stdout.strip())
    else:
        log(f"Command failed: {description}")
        if result.stderr.strip():
            print("STDERR:")
            print(result.stderr.strip())
        elif result.stdout.strip():
            # Some errors may go to stdout
            print("STDOUT (unexpected):")
            print(result.stdout.strip())

def extract_id(description: str, yq_path: str, *args) -> str:
    full_args = list(args) + ["-o", "yaml"]
    try:
        result = subprocess.run(full_args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=True, text=True)
        parsed_yaml = yaml.safe_load(result.stdout)

        if isinstance(parsed_yaml, list):
            first_item = parsed_yaml[0] if parsed_yaml else None
            id_val = first_item.get("id") if first_item else None
        elif isinstance(parsed_yaml, dict):
            if yq_path.startswith(".items"):
                items = parsed_yaml.get("items", [])
                id_val = items[0].get("id") if items else None
            elif yq_path.startswith("."):
                id_val = parsed_yaml.get(yq_path.lstrip("."))
            else:
                id_val = parsed_yaml.get(yq_path)
        else:
            id_val = None

        if not id_val:
            raise ValueError("ID not found")

        log(f"{description} ID: {id_val}")
        return id_val

    except Exception as e:
        log(f"{description} Failed to extract ID: {e}")
        sys.exit(1)

def run():
    global ROUTING_TABLE_ID, ROUTING_TABLE_ID_2, NETWORK_ID, ROUTE_ID

    run_command("Set project ID", "success", "./bin/stackit", "config", "set", "--project-id", PROJECT_ID)

    ROUTING_TABLE_ID = extract_id("Create routing-table rt_test", ".id",
        "./bin/stackit", "routing-table", "create", "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "--name", "rt_test", "-y")

    NETWORK_ID = extract_id("Create network with RT ID", ".id",
        "./bin/stackit", "network", "create", "--name", "network-rt", "--routing-table-id", ROUTING_TABLE_ID, "-y")

    run_command("List networks (check RT ID shown)", "success", "./bin/stackit", "network", "list", "-o", "pretty")
    run_command("Describe network", "success", "./bin/stackit", "network", "describe", NETWORK_ID)

    ROUTING_TABLE_ID_2 = extract_id("Create routing-table rt_test_2", ".id",
        "./bin/stackit", "routing-table", "create", "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "--name", "rt_test_2", "-y")

    run_command("Update network with RT 2 ID", "success",
        "./bin/stackit", "network", "update", NETWORK_ID, "--routing-table-id", ROUTING_TABLE_ID_2, "-y")

    run_command("Describe routing-table 1", "success",
        "./bin/stackit", "routing-table", "describe", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-o", "pretty")

    run_command("Describe routing-table 2", "success",
        "./bin/stackit", "routing-table", "describe", ROUTING_TABLE_ID_2,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-o", "pretty")

    run_command("List routing-tables", "success",
        "./bin/stackit", "routing-table", "list", "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "-o", "pretty")

    run_command("Delete second routing-table", "success",
        "./bin/stackit", "routing-table", "delete", ROUTING_TABLE_ID_2,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-y")

    run_command("Update RT: disable dynamic-routes", "success",
        "./bin/stackit", "routing-table", "update", ROUTING_TABLE_ID, "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "--description", "Test desc", "--non-dynamic-routes", "-y")

    run_command("Update RT: re-enable dynamic-routes", "success",
        "./bin/stackit", "routing-table", "update", ROUTING_TABLE_ID, "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "--description", "Test desc", "-y")

    run_command("Update RT: name", "success",
        "./bin/stackit", "routing-table", "update", ROUTING_TABLE_ID, "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "--name", "rt_test", "-y")

    run_command("Update RT: labels + name", "success",
        "./bin/stackit", "routing-table", "update", ROUTING_TABLE_ID, "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "--labels", "xxx=yyy,zzz=bbb", "--name", "rt_test", "-y")

    ROUTE_ID = extract_id("Create route with next-hop IPv4", ".items.0.id",
        "./bin/stackit", "routing-table", "route", "create", "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-y",
        "--destination-type", "cidrv4", "--destination-value", "0.0.0.0/0",
        "--nexthop-type", "ipv4", "--nexthop-value", "10.1.1.0")

    run_command("Create route with next-hop blackhole", "success",
        "./bin/stackit", "routing-table", "route", "create", "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-y",
        "--destination-type", "cidrv4", "--destination-value", "0.0.0.0/0", "--nexthop-type", "blackhole")

    run_command("Create route with next-hop internet", "success",
        "./bin/stackit", "routing-table", "route", "create", "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-y",
        "--destination-type", "cidrv4", "--destination-value", "0.0.0.0/0", "--nexthop-type", "internet")

    run_command("Negative test: invalid next-hop", "fail",
        "./bin/stackit", "routing-table", "route", "create", "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID,
        "--destination-type", "cidrv4", "--destination-value", "0.0.0.0/0", "--nexthop-type", "error")

    run_command("Negative test: invalid destination-type", "fail",
        "./bin/stackit", "routing-table", "route", "create", "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID,
        "--destination-type", "error", "--destination-value", "0.0.0.0/0", "--nexthop-type", "internet")

    run_command("List all routing-table routes", "success",
        "./bin/stackit", "routing-table", "route", "list", "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-o", "pretty")

    run_command("Describe route", "success",
        "./bin/stackit", "routing-table", "route", "describe", ROUTE_ID,
        "--routing-table-id", ROUTING_TABLE_ID, "--network-area-id", NETWORK_AREA_ID,
        "--organization-id", ORG_ID, "-o", "pretty")

    run_command("Update route labels", "success",
        "./bin/stackit", "routing-table", "route", "update", ROUTE_ID, "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID,
        "--labels", "key=value,foo=bar", "-y")

    run_command("Delete route", "success",
        "./bin/stackit", "routing-table", "route", "delete", ROUTE_ID, "--routing-table-id", ROUTING_TABLE_ID,
        "--network-area-id", NETWORK_AREA_ID, "--organization-id", ORG_ID, "-y")

    log("Cleanup: Removing all routing-tables named rt_test or rt_test_2.")
    cleanup_entities("routing-table", ["rt_test", "rt_test_2"],
                     ["--organization-id", ORG_ID, "--network-area-id", NETWORK_AREA_ID])

    log("Cleanup: Removing all networks named network-rt.")
    cleanup_entities("network", ["network-rt"], [])

    log("All tests finished successfully.")

def cleanup_entities(entity_type, name_list, extra_args):
    result = subprocess.run(["./bin/stackit", entity_type, "list", "-o", "yaml"] + extra_args,
                            stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    items = yaml.safe_load(result.stdout)
    for item in items:
        if item.get("name") in name_list:
            entity_id = item.get("id")
            cmd = ["./bin/stackit", entity_type, "delete", entity_id] + extra_args + ["-y"]
            run_command(f"Cleanup delete {entity_type} {item['name']}", "success", *cmd)

if __name__ == "__main__":
    run()
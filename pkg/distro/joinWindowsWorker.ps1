echo "join windows worker script";

$daemonJsonPath = "C:\ProgramData\Docker\config\daemon.json";

# Ensure the directory exists
if (-Not (Test-Path "C:\ProgramData\Docker\config")) {
    New-Item -ItemType Directory -Path "C:\ProgramData\Docker\config" -Force;
};

# Define the configuration
$daemonConfig = @{
    hosts = @("tcp://0.0.0.0:2375", "npipe://")
};

# Convert the configuration to JSON
$daemonConfigJson = $daemonConfig | ConvertTo-Json -Depth 3;

# Write the configuration to the file
$daemonConfigJson | Out-File -FilePath $daemonJsonPath -Encoding UTF8;

# Restart the Docker service to apply the new configuration
Restart-Service docker;
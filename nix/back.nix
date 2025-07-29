{
  pkgs,
  lib,
  self,
}: let
  appname = "eleakxir";
  backend = pkgs.buildGoModule {
    pname = appname + "-backend";
    version = "0.1.0";
    src = ../back;
    vendorHash = "sha256-xTVMERVvftdLJ8gT6OQRE+CrsHGOiYQ8tFPMTZc4A9U=";
    buildInputs = [
      pkgs.arrow-cpp
      pkgs.duckdb
    ];
  };
in {
  package = backend;

  nixosModule = {config, ...}: {
    options.services."${appname}-backend" = {
      enable = lib.mkEnableOption "Enable";
      port = lib.mkOption {
        type = lib.types.port;
        default = 8080;
        description = "Port on which the backend will listen.";
      };
      user = lib.mkOption {
        type = lib.types.str;
        default = appname + "-backend";
        description = "User under which the backend service will run.";
      };
      group = lib.mkOption {
        type = lib.types.str;
        default = appname + "-backend";
        description = "Group under which the backend service will run.";
      };
      leakPath = lib.mkOption {
        type = lib.types.str;
        default = "/var/lib/${appname}-backend/leaks";
        description = "Path to the leaks directory.";
      };
      cachePath = lib.mkOption {
        type = lib.types.str;
        default = "/var/lib/${appname}-backend/cache";
        description = "Path to the cache directory.";
      };
    };

    config = lib.mkIf config.services."${appname}-backend".enable {
      users.users."${config.services."${appname}-backend".user}" = {
        isSystemUser = true;
        group = config.services."${appname}-backend".group;
      };

      systemd.tmpfiles.rules = [
        "d /var/lib/${appname}-backend 0755 ${config.services."${appname}-backend".user} ${config.services."${appname}-backend".group} -"
        "d ${config.services."${appname}-backend".leakPath} 0775 ${config.services."${appname}-backend".user} ${config.services."${appname}-backend".group} -"
        "d ${config.services."${appname}-backend".cachePath} 0775 ${config.services."${appname}-backend".user} ${config.services."${appname}-backend".group} -"
      ];

      systemd.services."${appname}-backend" = {
        description = "${appname} Backend Service";
        after = ["network.target"];
        wantedBy = ["multi-user.target"];
        serviceConfig = {
          ExecStart = "${self.packages.${pkgs.system}."backend"}/bin/cmd";
          Restart = "always";
          User = config.services."${appname}-backend".user;
          Group = config.services."${appname}-backend".group;
          DynamicUser = false;
          StateDirectory = appname + "-backend";
          ReadWritePaths = ["/var/lib/${appname}-backend"];
          Environment = [
            "LEAK_DIRECTORY=${config.services."${appname}-backend".leakPath}"
            "CACHE_DIRECTORY=${config.services."${appname}-backend".cachePath}"
            "PORT=${toString config.services."${appname}-backend".port}"
          ];
        };
      };
    };
  };

  darwinModule = {
    config,
    lib,
    pkgs,
    ...
  }: {
    options.services."${appname}-backend" = {
      enable = lib.mkEnableOption "Enable";
      port = lib.mkOption {
        type = lib.types.port;
        default = 8080;
        description = "Port on which the backend will listen.";
      };
      leakPath = lib.mkOption {
        type = lib.types.str;
        default = "";
        description = "Path to the leaks directory.";
      };
      cachePath = lib.mkOption {
        type = lib.types.str;
        default = "";
        description = "Path to the cache directory.";
      };
    };

    config = lib.mkIf config.services."${appname}-backend".enable {
      launchd.user.agents."${appname}-backend" = {
        command = "${self.packages.${pkgs.system}."backend"}/bin/cmd";
        serviceConfig = {
          KeepAlive = true;
          EnvironmentVariables = {
            LEAK_DIRECTORY = config.services."${appname}-backend".leakPath;
            CACHE_DIRECTORY = config.services."${appname}-backend".cachePath;
            PORT = toString config.services."${appname}-backend".port;
          };
        };
      };
    };
  };
}

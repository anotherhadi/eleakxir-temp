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
    vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB=";
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
      systemd.services."${appname}-backend" = {
        description = "${appname} Backend Service";
        after = ["network.target"];
        wantedBy = ["multi-user.target"];
        serviceConfig = {
          ExecStart = "${self.packages.${pkgs.system}."backend"}/bin/cmd";
          Restart = "always";
          User = config.services."${appname}-backend".user;
          Group = config.services."${appname}-backend".group;
          DynamicUser = true;
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
}

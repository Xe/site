self:
{ config, lib, ... }:
with lib;
let cfg = config.xeserv.services.xesite;
in {
  options.xeserv.services.xesite = {
    enable = mkEnableOption "Activates my personal website";
    useACME = mkEnableOption "Enables ACME for cert stuff";

    port = mkOption {
      type = types.port;
      default = 32837;
      example = 9001;
      description = "The port number xesite should listen on for HTTP traffic";
    };

    domain = mkOption {
      type = types.str;
      default = "${config.networking.hostName}.shark-harmonic.ts.net";
      example = "xeiaso.net";
      description =
        "The domain name that nginx should check against for HTTP hostnames";
    };

    sockPath = mkOption rec {
      type = types.str;
      default = "/srv/within/run/xesite.sock";
      example = default;
      description = "The unix domain socket that xesite should listen on";
    };
  };

  config = mkIf cfg.enable {
    users.users.xesite = {
      createHome = true;
      description = "github.com/Xe/site";
      isSystemUser = true;
      group = "within";
      home = "/srv/within/xesite";
      extraGroups = [ "keys" ];
    };

    systemd.services.xesite = {
      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        User = "xesite";
        Group = "within";
        Restart = "on-failure";
        WorkingDirectory = "/srv/within/xesite";
        RestartSec = "30s";
        Type = "notify";

        # Security
        CapabilityBoundingSet = "";
        DeviceAllow = [ ];
        NoNewPrivileges = "true";
        ProtectControlGroups = "true";
        ProtectClock = "true";
        PrivateDevices = "true";
        PrivateUsers = "true";
        ProtectHome = "true";
        ProtectHostname = "true";
        ProtectKernelLogs = "true";
        ProtectKernelModules = "true";
        ProtectKernelTunables = "true";
        ProtectSystem = "true";
        ProtectProc = "invisible";
        RemoveIPC = "true";
        RestrictSUIDSGID = "true";
        RestrictRealtime = "true";
        SystemCallArchitectures = "native";
        SystemCallFilter = [
          "~@reboot"
          "~@module"
          "~@mount"
          "~@swap"
          "~@resources"
          "~@cpu-emulation"
          "~@obsolete"
          "~@debug"
          "~@privileged"
        ];
        UMask = "007";
      };

      script = let site = self.packages.${system}.default;
      in ''
        [ -f /srv/within/xesite/.env ] && export $(cat /srv/within/xesite/.env | xargs)
        export SOCKPATH=${cfg.sockPath}
        export DOMAIN=${toString cfg.domain}
        cd ${site}
        exec ${site}/bin/xesite
      '';
    };

    services.nginx.virtualHosts."xelaso.net" = let
      proxyOld = {
        proxyPass = "http://unix:${toString cfg.sockPath}";
        proxyWebsockets = true;
      };
    in {
      locations."/jsonfeed" = proxyOld;
      locations."/.within/health" = proxyOld;
      locations."/.within/website.within.xesite/new_post" = proxyOld;
      locations."/blog.rss" = proxyOld;
      locations."/blog.atom" = proxyOld;
      locations."/blog.json" = proxyOld;
      locations."/".extraConfig = ''
        return 301 https://xeiaso.net$request_uri;
      '';
      forceSSL = cfg.useACME;
      useACMEHost = "xeiaso.net";
      extraConfig = ''
        access_log /var/log/nginx/xesite_old.access.log;
      '';
    };

    services.nginx.virtualHosts."christine.website" = let
      proxyOld = {
        proxyPass = "http://unix:${toString cfg.sockPath}";
        proxyWebsockets = true;
      };
    in {
      locations."/jsonfeed" = proxyOld;
      locations."/.within/health" = proxyOld;
      locations."/.within/website.within.xesite/new_post" = proxyOld;
      locations."/blog.rss" = proxyOld;
      locations."/blog.atom" = proxyOld;
      locations."/blog.json" = proxyOld;
      locations."/".extraConfig = ''
        return 301 https://xeiaso.net$request_uri;
      '';
      forceSSL = cfg.useACME;
      useACMEHost = "christine.website";
      extraConfig = ''
        access_log /var/log/nginx/xesite_old.access.log;
      '';
    };

    services.nginx.virtualHosts."xeiaso.net" = {
      locations."/" = {
        proxyPass = "http://unix:${toString cfg.sockPath}";
        proxyWebsockets = true;
      };
      forceSSL = cfg.useACME;
      useACMEHost = "xeiaso.net";
      extraConfig = ''
        access_log /var/log/nginx/xesite.access.log;
      '';
    };
  };
}

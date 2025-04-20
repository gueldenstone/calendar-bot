{
  description = "Go development project template";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          name = "golang-devshell";

          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.golangci-lint
            pkgs.air
          ];

          shellHook = ''
            export GOPATH=$PWD/.gopath
            export PATH=$GOPATH/bin:$PATH
            mkdir -p .gopath/bin
            mkdir -p internal cmd pkg
            echo "🐹 Go devshell ready."
          '';
        };
      }
    );
}

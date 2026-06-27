{
  description = "MoneyRabbit dev shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go          # 1.26.x
            nodejs_24   # 24.x
            pnpm
            goose       # DB マイグレーション
            atlas       # ent スキーマ → SQL 差分生成
          ];

          shellHook = ''
            export GOPATH="$HOME/go"
            export PATH="$GOPATH/bin:$PATH"

            # swag は nixpkgs 未収録。初回のみ手動インストールが必要:
            #   go install github.com/swaggo/swag/cmd/swag@latest
            if ! command -v swag &>/dev/null; then
              echo "[moneyrabbit] swag not found. Run: go install github.com/swaggo/swag/cmd/swag@latest"
            fi
          '';
        };
      }
    );
}

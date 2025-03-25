{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  # buildInputs = with pkgs; [openssl pkg-config];
  #
  # RUST_BACKTRACE = 1;
}

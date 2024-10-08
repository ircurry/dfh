{
  version,
  lib,
  makeWrapper,
  buildGoModule,
  hyprland
}: buildGoModule {
  pname = "dfh";
  inherit version;

  src = ./.;

  nativeBuildInputs = [
    makeWrapper
  ];
  
  postFixup = ''
    wrapProgram $out/bin/dfh \
      --prefix PATH : "${
        lib.makeBinPath [
          hyprland
        ]
      }"
    wrapProgram $out/bin/hyprdock \
      --prefix PATH : "${
        lib.makeBinPath [
          hyprland
        ]
      }"
  '';

  vendorHash = null;
  
  meta = {
    description = "Helper for ircurry's dotfiles and desktop";
    homepage = "https://github.com/ircurry/dfh/";
    licence = lib.licences.mit;
  };
}

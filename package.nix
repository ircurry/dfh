{
  version,
  lib,
  makeWrapper,
  buildGoModule,
  hyprland,
  wlr-randr
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
          wlr-randr
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

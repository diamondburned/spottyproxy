{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
	buildInputs = with pkgs; [
		nodejs
		go
		gopls
	];

	shellHook = ''
		PATH="$PWD/node_modules/.bin:$PATH"
	'';
}

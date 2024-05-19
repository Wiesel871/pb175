bin_SCRIPT=bazos
noinst_DATA=pfp

install-exec-local:
	sh ./install-sh -m +x ./bazos $(prefix)/bazos

pfp:
	sh ./install-sh -d $(prefix)/images
	sh ./install-sh -d $(prefix)/images/0
	sh ./install-sh base_pfp.png $(prefix)/images/

install-data-local: pfp

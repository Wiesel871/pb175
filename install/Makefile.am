bin_SCRIPT=bazos
noinst_DATA=pfp admin

install-exec-local:
	sh ./install-sh -m +x ./bazos $(prefix)/bazos

pfp:
	sh ./install-sh -d $(prefix)/images
	sh ./install-sh -d $(prefix)/images/0
	sh ./install-sh -m 600 base_pfp.png $(prefix)/images/

admin:
	sh ./install-sh -m 600 admin.txt $(prefix)

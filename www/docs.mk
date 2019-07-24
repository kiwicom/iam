docs/publish:
	cd www && ./publish.sh

docs/serve:
	# -D shows draft pages as well
	cd www && hugo server -D

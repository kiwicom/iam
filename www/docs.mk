install_hugo:
	curl -sfL https://install.goreleaser.com/github.com/gohugoio/hugo.sh | sh

docs/publish:
	cd www && ./publish.sh

docs/serve:
	# -D shows draft pages as well
	cd www && hugo server -D
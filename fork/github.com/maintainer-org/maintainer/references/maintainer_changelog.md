## maintainer changelog

Generate CHANGELOG.md for your repository.

### Synopsis


changelog subcommand will generate CHANGELOG.md for your repository, it is supported
via github_changelog_generator, so you need to install it before the subcommand is called.

In the future, maintainer will support install this dependency automatically.

```
maintainer changelog
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.maintainer.yaml)
      --token string    The token in GitHub.To make more than 50 requests per hour the GitHub token is required.You can generate it at: https://github.com/settings/tokens/new.
```

### SEE ALSO
* [maintainer](maintainer.md)	 - Help you to be a qualified maintainer.

###### Auto generated by spf13/cobra on 23-May-2019
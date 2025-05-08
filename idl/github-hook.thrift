namespace go github_hook


service GithubHook {
    string ReleaseHook(1: string name) (api.get="/github-release");
}
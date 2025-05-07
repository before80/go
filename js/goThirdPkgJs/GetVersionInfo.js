() => {
    let versionE = document.querySelector(`span[data-test-id="UnitHeader-version"] > a`)
    let version = ""
    if (versionE) {
        if (versionE.querySelector(":scope > span")) {
            versionE.querySelector(":scope > span").remove()
        }
        version = versionE.textContent.trim()
    }
    let commitTimeE = document.querySelector(`span[data-test-id="UnitHeader-commitTime"]`)
    let commitTime = ""
    if (commitTimeE) {
        commitTime = commitTimeE.textContent.replace("Published:", "").trim()
    }

    let repoE = document.querySelector("div.UnitMeta-repo > a")
    let repo = ""
    if (repoE) {
        repo = repoE.href.trim()
    }

    return {
        version: version,
        commit_time: commitTime,
        repo:repo,
    }
}
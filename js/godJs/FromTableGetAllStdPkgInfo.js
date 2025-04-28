() => {
    // const btn = document.querySelector("button.js-showInternalDirectories")
    // if (btn) {
    //     if (btn.textContent.trim() === "Show internal") {
    //         btn.click()
    //         setTimeout(function (){}, 1000)
    //     }
    // }
    const table = document.querySelector("table.UnitDirectories-table")
    let pkgInfos = []
    let curTopPkgName = ""
    if (!table) {
        return
    }
    const trs = table.querySelectorAll(":scope > tbody > tr")
    if (trs.length === 0) {
        return
    }
    trs.forEach(tr => {
        console.log("run here 2")
        let desc = ""
        let pkgName = ""
        let url = ""
        let children = []
        const td1 = tr.querySelector(":scope > td:first-child")
        let isSub = false
        if (td1) {
            const div1 = td1.querySelector(":scope > div:first-child")
            if (div1) {
                if (div1.getAttribute("class") === "UnitDirectories-subdirectory") {
                    isSub = true
                    if (td1.querySelector("a")) {
                        pkgName = td1.querySelector("a").textContent.trim()
                        url = td1.querySelector("a").href.trim()
                    }
                } else {
                    if (td1.querySelector("a")) {
                        pkgName = td1.querySelector("a").textContent.trim()
                        url = td1.querySelector("a").href.trim()
                        curTopPkgName = pkgName
                    } else {
                        if (td1.querySelector("span")) {
                            pkgName = td1.querySelector("span").textContent.trim()
                            url = ""
                            curTopPkgName = pkgName
                        }
                    }
                }
            }
        }
        const td2 = tr.querySelector(":scope > td:nth-child(2)")
        if (td2) {
            desc = td2.textContent.trim()
        }
        if (pkgName) {
            if (isSub) {
                pkgInfos.push({
                    pkg_name: pkgName,
                    url: url,
                    desc: desc,
                    children: [],
                })
                for (let i = 0; i < pkgInfos.length; i++) {
                    if (pkgInfos[i].pkg_name === curTopPkgName) {
                        pkgInfos[i].children.push(pkgName);
                    }
                }
            } else {
                pkgInfos.push({
                    pkg_name: pkgName,
                    url: url,
                    desc: desc,
                    children: [],
                })
            }
        }
    })

    console.log("pkgInfos=", pkgInfos)
    return pkgInfos
}
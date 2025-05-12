() => {
    const curPageUrl = "%s"
    // const curPageUrl = "https://docs.python.org/zh-cn/3.13/index.html"
    let baseUrl = curPageUrl.replace(/\/index\.html$/, '');
    baseUrl = baseUrl.replace(/\/$/, '');
    baseUrl = baseUrl + '/'
    let menuInfos = []
    let exists = {}
    document.querySelectorAll("table.contentstable").forEach((t, i) => {
        const ps = t.querySelectorAll("p")
        if (ps.length > 0) {
            ps.forEach(p => {
                const a = p.querySelector('a')
                const menu_name = a.textContent.trim().replace(/[\"\'\/\\]/g,'')
                const urls = a.href.trim().split("#")
                const url = urls[0]
                let filename = url.replace(baseUrl, '')
                    .replace(/\/index\.html$/, '')
                    .replace(/\.html$/, '')
                    .replace(/[\.\/]/g, '_')

                if (!exists[url] && filename === "howto") {
                    if (i === 0 || (i > 0 && ["术语对照表", "Python 的历史与许可证"].includes(menu_name))) {
                        menuInfos.push({
                            menu_name: menu_name,
                            top_menu_name: menu_name,
                            is_top_menu: 1,
                            filename: filename,
                            weight: (i + 1) * 10,
                            dir: filename,
                            url: url,
                        })
                        exists[url] = true
                    }
                }
            })
        }
    })
    console.log(menuInfos)
    return menuInfos
}

() => {
    // toctree-wrapper compound
    const uls = document.querySelector("ul.simple")
    let menuInfos = []
    let exists = {}
    if (uls) {
        uls.querySelectorAll(":scope > li").forEach((li, i) => {
            const a = li.querySelector('a')
            const menu_name = a.textContent.trim().replace(/[\"\'\/\\]/g,'')
            const urls = a.href.trim().split("#")
            const url = urls[0]
            const names = url.split('/')
            let filename = names[names.length - 1].replace(/\.html$/, '')
                .replace(/[\.\/]/g, '_')
            const noJhaoUrl = url.split("#")[0]
            if (!exists[noJhaoUrl]) {
                menuInfos.push({
                    menu_name: menu_name,
                    is_top_menu: 2,
                    top_menu_name: "",
                    filename: filename,
                    weight: (i + 1) * 10,
                    dir: "",
                    url: url,
                })
                exists[noJhaoUrl] = true
            }
        })
    } else {
        //const curPageUrl = "%s"
    }
    console.log(menuInfos)
    return menuInfos
}


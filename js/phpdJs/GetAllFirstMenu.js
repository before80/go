() => {
    const lis = document.querySelectorAll("#index > ul.chunklist.chunklist_set > li")
    let menuInfos = []
    if (lis.length > 0) {
        lis.forEach((li, i) => {
            const a = li.querySelector(":scope > a")
            if (a) {
                const menuName = a.textContent.trim()
                const url = a.href.trim()
                let urls = url.split("/")
                const filename = urls[urls.length - 1].replace(/\.php$/, "").replace(/\./g,"_")
                // if (filename === "funcref") {
                //     menuInfos.push({
                //         menu_name: menuName,
                //         filename: filename,
                //         is_top_menu: 1,
                //         dir : "",
                //         url: url,
                //         weight: i + 1
                //     })
                // } else {
                //
                // }
                menuInfos.push({
                    menu_name: menuName,
                    filename: filename,
                    is_top_menu: 1,
                    dir : "",
                    url: url,
                    weight: i + 1
                })

            }
        })
    }
    console.log(menuInfos)
    return menuInfos
}
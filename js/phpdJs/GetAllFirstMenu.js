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

                menuInfos.push({
                    menu_name: menuName,
                    filename: filename,
                    url: url,
                    index: i + 1
                })
            }
        })
    }
    console.log(menuInfos)
    return menuInfos
}
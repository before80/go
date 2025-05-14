() => {
    let menuInfos = []
    const h1 = document.querySelector('h1')
    if (h1) {
        const lis = h1.parentElement.querySelectorAll(":scope > ul.chunklist > li")
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
                        is_top_menu: 2,
                        dir: "",
                        url: url,
                        weight: i + 1
                    })
                }
            })
        }
    }
    console.log(menuInfos)
    return menuInfos
}
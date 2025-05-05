() => {
    const as = document.querySelectorAll("nav.doctoc > ul > li a.expandable")
    if (as.length > 0) {
        as.forEach(a => {
            a.click()
        })
    }
}
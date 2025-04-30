() => {
    const layout = document.querySelector("#layout")
    let content = ""
    if (layout) {
        content = layout.textContent.trim()
    }
    return content
}
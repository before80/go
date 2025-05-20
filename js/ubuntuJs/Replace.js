
function removeSomething() {
    const header = document.querySelector("header")
    if (header) {
        header.remove()
    }
    const toc = document.querySelector("#toc")
    if (toc) {
        toc.remove()
    }
    const footer = document.querySelector("footer")
    if (footer) {
        footer.remove()
    }


}

function replaceSomething() {
    const h4s = document.querySelectorAll("h4");
    if (h4s.length > 0) {
        h4s.forEach(h4 => {
            const h2 = document.createElement("h2")
            h2.innerHTML = h4.innerHTML
            h4.insertAdjacentElement("beforebegin", h2)
            h4.remove()
        })
    }
    const pres = document.querySelectorAll("h2 + pre");
    if (pres.length > 0) {
        pres.forEach(pre => {
            pre.innerHTML = pre.innerHTML.replace(/([^。])\n\s+/g, `$1`)
            pre.childNodes.forEach(originNode => {
                if (originNode.nodeType === Node.TEXT_NODE) {

                    console.log("node.nodeValue=",`"${originNode.nodeValue}"`)
                    if (originNode.nodeValue.includes("       ")) {
                        originNode.nodeValue = originNode.nodeValue.replace(/       /g,"")
                    }
                    if (originNode.nodeValue.includes("'")) {
                        originNode.nodeValue = originNode.nodeValue.replace(/'/g,"`")
                    }

                    if (originNode.nodeValue.includes('\n')) {
                        // 分割文本并创建新节点数组
                        const parts = originNode.nodeValue.split('\n');
                        const newNodes = [];

                        parts.forEach((part, index) => {
                            // 添加文本部分
                            newNodes.push(document.createTextNode(part));

                            // 除了最后一部分，为每部分后面添加 <br> 标签
                            if (index < parts.length - 1) {
                                newNodes.push(document.createElement('br'));
                                newNodes.push(document.createElement('br'));
                            }
                        });

                        // 用新节点替换原始文本节点
                        const parent = originNode.parentNode;
                        newNodes.forEach(node => parent.insertBefore(node, originNode));
                        parent.removeChild(originNode);
                    }
                }
            })

            const bs = pre.querySelectorAll("b")
            if (bs.length > 0) {
                bs.forEach(b => {
                    if (!b.querySelector("a")) {
                        const u = b.querySelector(":scope + u")
                        const span = document.createElement("span")
                        span.innerHTML = "\u0060" + b.innerHTML + "\u0060";
                        b.insertAdjacentElement("beforebegin", span)
                        b.remove()
                        // b.innerHTML = "\u0060" + b.innerHTML + "\u0060";
                    }
                })
            }

            const div = document.createElement("div")
            div.innerHTML = pre.innerHTML
            pre.insertAdjacentElement("beforebegin", div)
            pre.remove()
        })
    }
}

removeSomething();
replaceSomething();
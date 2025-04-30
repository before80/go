function removeSomething() {
    const nav = document.querySelector("nav.navbar.navbar-fixed-top")
    if (nav) {
        nav.remove()
    }

    const breadcrumbs = document.querySelector("#breadcrumbs")
    if (breadcrumbs) {
        breadcrumbs.remove()
    }

    const trick = document.querySelector("#trick")
    if (trick) {
        trick.remove()
    }

    const goto = document.querySelector("#goto")
    if (goto) {
        goto.remove()
    }

    const footer = document.querySelector("footer")
    if (footer) {
        footer.remove()
    }

    const layoutMenu = document.querySelector("#layout aside.layout-menu")
    if (layoutMenu) {
        layoutMenu.remove()
    }

    const contribute = document.querySelector("#layout-content div.contribute")
    if (contribute) {
        contribute.remove()
    }

    const userNotes = document.querySelector("#usernotes")
    if (userNotes) {
        userNotes.remove()
    }

    const changeLanguage = document.querySelector(".change-language")
    if (changeLanguage) {
        changeLanguage.remove()
    }

    const toTop = document.querySelector("#toTop")
    if (toTop) {
        toTop.remove()
    }

    const h1 = document.querySelector("h1")
    if (h1) {
        const chunklist = h1.parentElement.querySelector("ul.chunklist")
        if (chunklist) {
            chunklist.remove()
        }
        h1.remove()
    }


    const h2s = document.querySelectorAll("h2")
    if (h2s.length > 0) {
        h2s.forEach(h2 => {
            if (h2.textContent.startsWith('目录')) {
                h2.remove()
            }
        })
    }


    const btns = document.querySelectorAll("button")
    if (btns.length > 0) {
        btns.forEach(btn => {
            btn.remove()
        })
    }
    const layout = document.querySelector("#layout")
    if (layout) {
        if(layout.textContent.trim() === "") {
            console.log("无内容")
            // layout.insertAdjacentText("afterbegin","无内容")
        }
    }
}


function convertCodeToHtml(code, lang) {
    // 将 \t 替换为四个空格
    let converted = code.replace(/\t/g, '&nbsp;&nbsp;&nbsp;&nbsp;');
    // 将 \n 替换为 <br> 标签
    converted = converted.replace(/\n\)$/, '<br><br>)');
    converted = converted.replace(/\n\}$/, '<br><br>}');
    converted = converted.replace(/\n/g, '<br>');
    // 使用 <pre> 和 <code> 标签包裹代码，并添加语言类名
    const html = `<div><pre><code class="text-sm text-gray-800 bg-gray-200 p-4 rounded-md language-${lang}">${converted}</code></pre></div>`;
    return html;
}

function replaceDivClassSynopsis() {
    const divs = document.querySelectorAll("div.classsynopsis")
    if (divs.length > 0) {
        divs.forEach(div => {
            const html = convertCodeToHtml(div.textContent.trim(),"php")
            div.insertAdjacentHTML("afterend",html)
            div.remove()
        })
    }
}

function replaceShellCode() {
    const shellCodes = document.querySelectorAll("div.shellcode")
    if (shellCodes.length > 0) {
        shellCodes.forEach(sc => {
            const pre = sc.querySelector(":scope > pre")
            const div = document.createElement("div")
            const html = `<pre><code class="text-sm text-gray-800 bg-gray-200 p-4 rounded-md language-sh">${pre.innerHTML}</code></pre>`;
            div.insertAdjacentHTML("afterbegin", html)
            pre.insertAdjacentElement("afterend", div)
            pre.remove()
        })
    }
}

function replaceHtmlCode() {
    const htmlCodes = document.querySelectorAll("div.htmlcode")
    if (htmlCodes.length > 0) {
        htmlCodes.forEach(hc => {
            const pre = hc.querySelector(":scope > pre")
            const div = document.createElement("div")
            const html = `<pre><code class="text-sm text-gray-800 bg-gray-200 p-4 rounded-md language-html">${pre.innerHTML}</code></pre>`;
            div.insertAdjacentHTML("afterbegin", html)
            pre.insertAdjacentElement("afterend", div)
            pre.remove()
        })
    }
}

function replacePHPCode() {
    const phpCodes = document.querySelectorAll("div.phpcode")
    if (phpCodes.length > 0) {
        phpCodes.forEach(pc => {
            const code = pc.querySelector(":scope > code")
            const div = document.createElement("div")
            const html = `<pre><code class="text-sm text-gray-800 bg-gray-200 p-4 rounded-md language-php">${code.innerHTML}</code></pre>`;
            div.insertAdjacentHTML("afterbegin", html)
            code.insertAdjacentElement("afterend", div)
            code.remove()
        })
    }
}

function replaceCData() {
    const cds = document.querySelectorAll("div.cdata")
    if (cds.length > 0) {
        cds.forEach(cd => {
            const pre = cd.querySelector(":scope > pre")
            const div = document.createElement("div")
            const html = `<pre><code class="text-sm text-gray-800 bg-gray-200 p-4 rounded-md language-txt">${pre.innerHTML}</code></pre>`;
            div.insertAdjacentHTML("afterbegin", html)
            pre.insertAdjacentElement("afterend", div)
            pre.remove()
        })
    }
}

function addBlockquote() {
    const divs = document.querySelectorAll("div.tip,div.warning")
    if (divs.length > 0) {
        divs.forEach(div => {
            const blockquote = document.createElement('blockquote')
            const newDiv = document.createElement('div')
            newDiv.insertAdjacentHTML("afterbegin", div.innerHTML)
            blockquote.appendChild(newDiv)
            div.insertAdjacentElement("afterend", blockquote)
            div.remove()
        })
    }
}

function replacePContent() {
    const pTags = document.getElementsByTagName('p');
    for (let i = 0; i < pTags.length; i++) {
        const p = pTags[i];
        for (let j = 0; j < p.childNodes.length; j++) {
            const node = p.childNodes[j];
            if (node.nodeType === Node.TEXT_NODE) {
                node.nodeValue = node.nodeValue.replace(/(_[_a-zA-Z:\(\)\?\$]+)/g, '`$1`');
                // node.nodeValue = node.nodeValue.replace(/(__construct\(\))/g, '`$1`');
            }
        }
    }
}

// 在p标签前面添加一个span，其内容为&zeroWidthSpace;用于后期在markdown中替换成Tab符号
function replaceP() {
    // 在p标签前面插入&zeroWidthSpace;
    document.querySelectorAll('p.simpara,p.para').forEach(function (p) {
        if (!['LI', "BLOCKQUOTE", "TH", "TD"].includes(p.parentElement.tagName)) {
            if (p.textContent.trim() !== "") {
                let newSpan = document.createElement('span');
                newSpan.textContent = '&zeroWidthSpace;';
                if (p.firstChild) {
                    p.insertBefore(newSpan, p.firstChild);
                } else {
                    // 如果 p 元素没有子节点，直接将新 span 元素添加到 p 元素中
                    p.appendChild(newSpan);
                }
            } else {
                p.remove()
            }
        }
    });
}

function replaceVar() {
    document.querySelectorAll("var").forEach(v => {
        v.innerHTML = "\u0060" + v.innerHTML + "\u0060";
    })
}


function addHeaderAnchorAndRemoveGenanchor() {
    const hs = document.querySelectorAll("h2,h3,h4,h5,h6")

    hs.forEach(h => {
        const genanchor = h.querySelector("a.genanchor")
        if (genanchor) {
            const link = genanchor.href
            //去除后面的 ¶
            const anchor = link.split("#")[1].trim()
            genanchor.remove()
            h.insertAdjacentText("beforeend", `{#${anchor}}`)
        }
    })
}


removeSomething();
replaceDivClassSynopsis();
replaceShellCode();
replaceHtmlCode();
replacePHPCode();
addBlockquote();
replacePContent();
replaceP();
replaceVar();
addHeaderAnchorAndRemoveGenanchor();
replaceCData();
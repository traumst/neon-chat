function escFocus(event) {
    if (event.key === 'Escape') {
        event.target.blur();
    } else {
        //console.log('escFocus key', event.key)
    }
}

function enterSubmit(event) {
    if (event.key === 'Enter') {
        event.target.submit();
    } else {
        //console.log('enterSubmit key', event.key)
    }
}

function removeElements({classSelector, idSelectorException}) {
    //console.log('removeElements', classSelector, "except", idSelectorException)
    const removeTargets = document.querySelectorAll(classSelector);
    for (const removeTarget of removeTargets) {
        if (!idSelectorException || removeTarget.id !== idSelectorException) {
            removeTarget.remove();
        }
    }
}

async function scrollIntoView(elemId) {
    const target = document.querySelectorAll(elemId);
    if (target.length != 1) {
        console.error(
            `scrollIntoView: exactly 1 element should have matched ${elemId},`, 
            `but found ${target.length}:`, target)
    } 
    if (target[0]) {
        // console.log(`scrollIntoView: scrolling element ${elemId} into view`)
        await target[0].scrollIntoView({ 
            behavior: 'smooth',
            block: 'center'
        });
        // add flash highlight
        setTimeout(() => {
            target[0].classList.add('flash-golden');
        }, 150); // so it does not flush offscreen while scrolling
        // remove flash highlight
        setTimeout(() => {
            target[0].classList.remove('flash-golden');
        }, 4000); // must be longer than animation duration
    }
}

function scrollToFirstChild(parent){
    if (!parent || !parent.children.length) {
        return;
    }
    const firstChild = parent.children[0];
    if (firstChild) {
        // console.log(`scrollToFirstChild: scrolling child ${firstChild.id} into view`);
        firstChild.scrollIntoView({ 
            behavior: 'smooth',
            block: 'start'
        });
    } else {
        console.error(`scrollToFirstChild: parent ${parent.id} not found`)
    }
}

function scrollToLastChild(parent) {
    if (!parent || !parent.children.length) {
        return;
    }
    const lastChild = parent.children[parent.children.length-1];
    if (lastChild) {
        // console.log('scrollToLastChild: scrolling to child', lastChild.id);
        lastChild.scrollIntoView({
            behavior: 'smooth',
            block: 'end'
        });
    } else {
        console.error('scrollToLastChild: not found parent', parent.id);
    }
}

function showScopedPopup(targetUnderlayId, targetPopupId) {
    console.log('showScopedPopup:', targetUnderlayId, targetPopupId);
    const userInfoUnderlay = document.getElementById(targetUnderlayId);
    while (userInfoUnderlay && userInfoUnderlay.classList.contains('hidden')) {
        userInfoUnderlay.classList.remove('hidden');
    }
    const userInfoPopup = document.getElementById(targetPopupId);
    while (userInfoPopup && userInfoPopup.classList.contains('hidden')) {
        userInfoPopup.classList.remove('hidden');
    }
}

function hideScopedPopup(targetUnderlayId, targetPopupId) {
    console.log('hideScopedPopup:', targetUnderlayId, targetPopupId);
    const userInfoPopup = document.getElementById(targetPopupId);
    if (userInfoPopup && !userInfoPopup.classList.contains('hidden')) {
        userInfoPopup.classList.add('hidden');
    }
    const userInfoUnderlay = document.getElementById(targetUnderlayId);
    if (userInfoUnderlay && !userInfoUnderlay.classList.contains('hidden')) {
        userInfoUnderlay.classList.add('hidden');
    }
}
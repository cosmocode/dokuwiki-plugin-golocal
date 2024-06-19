class GoLocal {

    /** @type {HTMLDialogElement} */
    #dialog;

    /** @type {URL} */
    #targetURL;

    /**
     *
     * @param {Event} event
     */
    constructor(event) {
        // decide if a dialog should be shown
        if (localStorage.getItem('golocal')) {
            // continue without dialog
            return;
        }

        // still here? continue with dialog
        this.#targetURL = event.target.href;
        event.preventDefault();

        // create HTML and attach to document
        this.#dialog = this.createDialog();
        const parent = document.body.querySelector('.dokuwiki') || document.body;
        parent.appendChild(this.#dialog);


        this.#dialog.showModal();
    }

    /**
     * Create the dialog HTML and attach event listeners
     *
     * @returns {HTMLDialogElement}
     */
    createDialog() {
        const dialog = document.createElement('dialog');
        dialog.className = 'golocal-dialog';

        const main = document.createElement('main');
        dialog.appendChild(main);

        const intro = document.createElement('div');
        intro.innerHTML = LANG.plugins.golocal.dialog_intro;
        main.appendChild(intro);

        const url = new URL(window.location.href);
        url.searchParams.set('do', 'golocal');
        const install = document.createElement('a');
        install.href = url.toString();
        install.className = 'install';
        install.textContent = LANG.plugins.golocal.dialog_install;
        main.appendChild(install);

        const footer = document.createElement('footer');
        dialog.appendChild(footer);

        const rememberLabel = document.createElement('label');

        const rememberCheckbox = document.createElement('input');
        rememberCheckbox.type = 'checkbox';
        rememberLabel.appendChild(rememberCheckbox);

        const rememberSpan = document.createElement('span');
        rememberSpan.textContent = LANG.plugins.golocal.dialog_remember;
        rememberLabel.appendChild(rememberSpan);

        footer.appendChild(rememberLabel);

        const continueButton = document.createElement('button');
        continueButton.textContent = LANG.plugins.golocal.dialog_continue;
        continueButton.addEventListener('click', (ev) => this.onContinue(ev, rememberCheckbox));
        continueButton.autofocus = true;
        footer.appendChild(continueButton);

        dialog.addEventListener('close', this.onClose.bind(this));
        return dialog;
    }

    /**
     * Continue to the target URL
     * @param {Event} ev The click event
     * @param {HTMLInputElement} rememberCheckbox The checkbox to remember the choice (if available)
     */
    onContinue(ev, rememberCheckbox = null) {
        if (rememberCheckbox && rememberCheckbox.checked) {
            localStorage.setItem('golocal', 'true');
        }

        this.#dialog.close();
        window.location.href = this.#targetURL;
    }


    onClose(e) {
        this.#dialog.remove();
    }
}

jQuery(function () {
    jQuery('a.windows').each(function () {
        // fix up all standard windows share links
        let $this = jQuery(this);
        $this.attr('href', $this.attr('href').replace('file:///', 'golocal://'));
    }).click(function (e) {
        new GoLocal(e);
    });
});

<?php

use dokuwiki\Extension\ActionPlugin;
use dokuwiki\Extension\Event;
use dokuwiki\Extension\EventHandler;

/**
 * DokuWiki Plugin golocal (Action Component)
 *
 * @license GPL 2 http://www.gnu.org/licenses/gpl-2.0.html
 * @author Andreas Gohr <dokuwiki@cosmocode.de>
 */
class action_plugin_golocal extends ActionPlugin
{
    /** @inheritDoc */
    public function register(EventHandler $controller)
    {
        $controller->register_hook('INIT_LANG_LOAD', 'AFTER', $this, 'handleInitLangLoad');

        $controller->register_hook('ACTION_ACT_PREPROCESS', 'BEFORE', $this, 'handleActionActPreprocess');
        $controller->register_hook('TPL_ACT_UNKNOWN', 'BEFORE', $this, 'handleTplActUnknown');
    }

    /**
     * Event handler for INIT_LANG_LOAD
     *
     * @see https://www.dokuwiki.org/devel:events:init_lang_load
     * @param Event $event Event object
     * @param mixed $param optional parameter passed when event was registered
     * @return void
     */
    public function handleInitLangLoad(Event $event, $param)
    {
        global $lang;
        $lang['js']['nosmblinks'] = '';
    }

    /**
     * Event handler for ACTION_ACT_PREPROCESS
     *
     * @see https://www.dokuwiki.org/devel:events:action_act_preprocess
     * @param Event $event Event object
     * @param mixed $param optional parameter passed when event was registered
     * @return void
     */
    public function handleActionActPreprocess(Event $event, $param)
    {
        if ($event->data !== 'golocal') return;
        $event->preventDefault();
        $event->stopPropagation();
    }

    public function handleTplActUnknown(Event $event, $param)
    {
        if ($event->data !== 'golocal') return;
        $event->preventDefault();
        $event->stopPropagation();

        $output = $this->locale_xhtml('download');
        $output = str_replace('DOWNLOADSHERE', $this->getDownloadLinks(), $output);
        echo $output;
    }

    /**
     * Build the list of downloads
     *
     * @todo this could be a syntax component
     * @return string
     */
    protected function getDownloadLinks()
    {

        $files = glob(__DIR__ . '/dl/*');

        $html = '<ul class="golocal-download">';
        foreach ($files as $file) {
            $file = basename($file);
            $url = DOKU_BASE . 'lib/plugins/golocal/dl/' . $file;
            [$os, $arch] = sexplode('_', (sexplode('-', basename($file, '.exe'), 2))[1], 2);

            $classes = join(' ', ['li', 'os-' . $os, 'arch-' . $arch]);

            $html .= '<li><div class="' . $classes . '">';
            $html .= inlineSVG(__DIR__ . '/icons/' . $os . '.svg');
            $html .= inlineSVG(__DIR__ . '/icons/' . $arch . '.svg');
            $html .= '<a href="' . $url . '">' . $file . '</a>';
            $html .= '</div></li>';

        }
        $html .= '</ul>';


        return $html;
    }


}

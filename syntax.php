<?php

use dokuwiki\Extension\SyntaxPlugin;

/**
 * DokuWiki Plugin golocal (Syntax Component)
 *
 * @license GPL 2 http://www.gnu.org/licenses/gpl-2.0.html
 * @author Andreas Gohr <dokuwiki@cosmocode.de>
 */
class syntax_plugin_golocal extends SyntaxPlugin
{
    /** @inheritDoc */
    public function getType()
    {
        return 'substition';
    }

    /** @inheritDoc */
    public function getPType()
    {
        return 'normal';
    }

    /** @inheritDoc */
    public function getSort()
    {
        return 150;
    }

    /** @inheritDoc */
    public function connectTo($mode)
    {
        $this->Lexer->addSpecialPattern('\\[\\[[C-Z]:\\\\[^]]*\\]\\]', $mode, 'plugin_golocal');
    }


    /** @inheritDoc */
    public function handle($match, $state, $pos, Doku_Handler $handler)
    {
        $match = substr($match, 2, -2);
        [$path, $title] = sexplode('|', $match, 2);
        $path = trim($path);
        $title = trim($title);
        if (!$title) $title = $path;

        return [$path, $title];
    }

    /** @inheritDoc */
    public function render($mode, Doku_Renderer $renderer, $data)
    {
        if ($mode == 'xhtml') {
            $params = [
                'href' => 'file:////' . str_replace(':', '/', str_replace('\\', '/', $data[0])),
                'title' => $data[0],
                'class' => 'windows'
            ];

            $renderer->doc .= '<a ' . buildAttributes($params) . '>' . hsc($data[1]) . '</a>';
        } else {
            $renderer->cdata($data[1]);
        }

        return true;
    }
}

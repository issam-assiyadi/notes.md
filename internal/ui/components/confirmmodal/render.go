package confirmmodal

import "github.com/awesome-gocui/gocui"

func (m *Modal) Render(g *gocui.Gui) error {
	if !m.open {
		return nil
	}

	m.confirm.Render(g, m.focus == fieldConfirm)
	m.dismiss.Render(g, m.focus == fieldDismiss)

	if v, err := g.View(m.messageViewName()); err == nil {
		v.Clear()
		_, _ = v.Write([]byte(m.message))
	}
	return nil
}

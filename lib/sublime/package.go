package sublime

type pkg struct {
	text.HasSettings
	keys.HasKeyBindings
	platformSet *text.HasSettings
	defaultSet  *text.HasSettings
	defaultKB   *keys.HasKeyBindings
	plugins     []plugin
}

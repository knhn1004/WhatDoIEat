package services

func StartServices() {
	InitOpenAI()
	InitCohere()
	InitSupabase()
}

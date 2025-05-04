package database

import (
	"time"

	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var threads = []model.Thread{
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "This is thread #1. Welcome to the discussion!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 1),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #2. Does anyone have tips for beginners?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 2),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #3. Looking for collaborators on a project.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 3),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #4. Weekly update on community progress.",
			ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 4),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #5. AMA (Ask Me Anything) starting now!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 5),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #6. Is this feature working for you?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 6),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #7. Survey results are out!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 7),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #8. Share your weekend projects!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 8),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #9. Major announcement coming soon.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 9),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #10. Let's talk about mental health.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 10),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #11. Happy to join this platform!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 11),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #12. Daily inspiration thread 🌟",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 12),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #13. Need help with deployment.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 13),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #14. Monthly recap and plans ahead.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 14),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #15. Advert: Get 20% off your first order!",
			ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 15),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "This is thread #1. Welcome to the discussion!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 1),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #2. Does anyone have tips for beginners?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 2),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #3. Looking for collaborators on a project.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 3),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #4. Weekly update on community progress.",
			ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 4),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #5. AMA (Ask Me Anything) starting now!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 5),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #6. Is this feature working for you?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 6),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #7. Survey results are out!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 7),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #8. Share your weekend projects!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 8),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #9. Major announcement coming soon.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 9),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #10. Let's talk about mental health.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 10),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #11. Happy to join this platform!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 11),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #12. Daily inspiration thread 🌟",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 12),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #13. Need help with deployment.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 13),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #14. Monthly recap and plans ahead.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 14),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #15. Advert: Get 20% off your first order!",
			ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 15),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "This is thread #1. Welcome to the discussion!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 1),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #2. Does anyone have tips for beginners?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 2),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #3. Looking for collaborators on a project.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 3),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #4. Weekly update on community progress.",
			ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 4),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #5. AMA (Ask Me Anything) starting now!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 5),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #6. Is this feature working for you?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 6),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #7. Survey results are out!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 7),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #8. Share your weekend projects!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 8),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #9. Major announcement coming soon.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 9),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #10. Let's talk about mental health.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 10),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #11. Happy to join this platform!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 11),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #12. Daily inspiration thread 🌟",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 12),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #13. Need help with deployment.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 13),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #14. Monthly recap and plans ahead.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 14),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #15. Advert: Get 20% off your first order!",
			ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 15),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "This is thread #1. Welcome to the discussion!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 1),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #2. Does anyone have tips for beginners?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 2),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #3. Looking for collaborators on a project.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 3),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #4. Weekly update on community progress.",
			ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 4),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #5. AMA (Ask Me Anything) starting now!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 5),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #6. Is this feature working for you?",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 6),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #7. Survey results are out!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 7),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #8. Share your weekend projects!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 8),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #9. Major announcement coming soon.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 9),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #10. Let's talk about mental health.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 10),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #11. Happy to join this platform!",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 11),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #12. Daily inspiration thread 🌟",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 12),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #13. Need help with deployment.",
			ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 13),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #14. Monthly recap and plans ahead.",
			ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 14),
		},
		{
			UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
			Content: "Thread #15. Advert: Get 20% off your first order!",
			ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
			CreatedAt: time.Now().Add(-time.Hour * 15),
		},
	}
	for _, element := range threads {
		db.Create(&element)
	}

	var categories = []model.ThreadCategory{
		{
			Name: "Anime",
		},
		{
			Name: "Manga",
		},
		{
			Name: "Manhwa",
		},
	}
	for _, element := range categories {
		db.Create(&element)
	}

	var communities = []model.Community{
		{
			Name: "Hololive",
			Description: "A very special place for everyone",
			Logo: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//hoshimachi-suisei-hololive.gif",
			Category: "Game",
			Banner: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(125).png",
			Rule: "Please support your oshis wholeheartedly",
		},
		{
			Name: "Hoyoverse",
			Description: "A place that talks about the best game creator",
			Logo: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(2).png",
			Category: "Game",
			Banner: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(3).png",
			Rule: "Enjoy hoyoverse games",
		},
	}
	for _, element := range communities {
		db.Create(&element)
	}
}
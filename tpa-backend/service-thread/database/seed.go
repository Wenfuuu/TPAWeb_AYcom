package database

import (
	// "time"

	"math/rand"
	"strconv"
	"time"

	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

func generateRandomCommunities() []model.Community {
    rand.Seed(time.Now().UnixNano())

    logos := []string{
        "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//hoshimachi-suisei-hololive.gif",
        "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(2).png",
    }

    banners := []string{
        "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(125).png",
        "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(3).png",
    }

    categoryIDs := []uuid.UUID{
        uuid.MustParse("928c8cc1-a75c-4fd2-a080-f2e4e090e698"),
        uuid.MustParse("b49f4771-6207-46a2-84ef-8208775ae922"),
    }

    communities := []model.Community{}
    for i := 0; i < 30; i++ {
        community := model.Community{
            Name:        	"Community " + strconv.Itoa(i+1), // Convert i+1 to string
            Description: 	"Description for community " + strconv.Itoa(i+1),
            Logo:        	logos[rand.Intn(len(logos))],
            Banner:      	banners[rand.Intn(len(banners))],
            CategoryID:  	categoryIDs[rand.Intn(len(categoryIDs))],
            Rule:        	"Random rule for community",
			CreatorID: 		uuid.New(),
        }
        communities = append(communities, community)
    }
    return communities
}

func Seed(db *gorm.DB) {
	// var threads = []model.Thread{
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "This is thread #1. Welcome to the discussion!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 1),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #2. Does anyone have tips for beginners?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 2),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #3. Looking for collaborators on a project.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 3),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #4. Weekly update on community progress.",
	// 		ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 4),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #5. AMA (Ask Me Anything) starting now!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 5),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #6. Is this feature working for you?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 6),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #7. Survey results are out!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 7),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #8. Share your weekend projects!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 8),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #9. Major announcement coming soon.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 9),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #10. Let's talk about mental health.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 10),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #11. Happy to join this platform!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 11),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #12. Daily inspiration thread 🌟",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 12),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #13. Need help with deployment.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 13),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #14. Monthly recap and plans ahead.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 14),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #15. Advert: Get 20% off your first order!",
	// 		ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 15),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "This is thread #1. Welcome to the discussion!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 1),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #2. Does anyone have tips for beginners?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 2),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #3. Looking for collaborators on a project.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 3),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #4. Weekly update on community progress.",
	// 		ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 4),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #5. AMA (Ask Me Anything) starting now!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 5),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #6. Is this feature working for you?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 6),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #7. Survey results are out!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 7),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #8. Share your weekend projects!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 8),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #9. Major announcement coming soon.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 9),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #10. Let's talk about mental health.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 10),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #11. Happy to join this platform!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 11),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #12. Daily inspiration thread 🌟",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 12),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #13. Need help with deployment.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 13),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #14. Monthly recap and plans ahead.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 14),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #15. Advert: Get 20% off your first order!",
	// 		ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 15),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "This is thread #1. Welcome to the discussion!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 1),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #2. Does anyone have tips for beginners?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 2),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #3. Looking for collaborators on a project.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 3),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #4. Weekly update on community progress.",
	// 		ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 4),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #5. AMA (Ask Me Anything) starting now!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 5),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #6. Is this feature working for you?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 6),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #7. Survey results are out!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 7),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #8. Share your weekend projects!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 8),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #9. Major announcement coming soon.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 9),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #10. Let's talk about mental health.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 10),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #11. Happy to join this platform!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 11),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #12. Daily inspiration thread 🌟",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 12),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #13. Need help with deployment.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 13),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #14. Monthly recap and plans ahead.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 14),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #15. Advert: Get 20% off your first order!",
	// 		ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 15),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "This is thread #1. Welcome to the discussion!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 1),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #2. Does anyone have tips for beginners?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 2),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #3. Looking for collaborators on a project.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 3),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #4. Weekly update on community progress.",
	// 		ReplyPermission: "admins", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 4),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #5. AMA (Ask Me Anything) starting now!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 5),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #6. Is this feature working for you?",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 6),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #7. Survey results are out!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 7),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #8. Share your weekend projects!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 8),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #9. Major announcement coming soon.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 9),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #10. Let's talk about mental health.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 10),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #11. Happy to join this platform!",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 11),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #12. Daily inspiration thread 🌟",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 12),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #13. Need help with deployment.",
	// 		ReplyPermission: "everyone", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 13),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #14. Monthly recap and plans ahead.",
	// 		ReplyPermission: "followers", IsAdvertisement: false, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 14),
	// 	},
	// 	{
	// 		UserID: uuid.New(), CategoryID: uuid.New(), CommunityID: uuid.New(),
	// 		Content: "Thread #15. Advert: Get 20% off your first order!",
	// 		ReplyPermission: "everyone", IsAdvertisement: true, Status: "published",
	// 		CreatedAt: time.Now().Add(-time.Hour * 15),
	// 	},
	// }
	// for _, element := range threads {
	// 	db.Create(&element)
	// }

	// var categories = []model.ThreadCategory{
		// {
		// 	Name: "Anime",
		// },
		// {
		// 	Name: "Manga",
		// },
	// 	{
	// 		Name: "Manhwa",
	// 	},
	// }
	// for _, element := range categories {
	// 	db.Create(&element)
	// }

	// var communityCategories = []model.CommunityCategory{
	// 	{
	// 		Name: "VTuber",
	// 	},
	// 	{
	// 		Name: "Game",
	// 	},
	// }
	// for _, element := range communityCategories {
	// 	db.Create(&element)
	// }

	var communities = []model.Community{
		{
			Name: "Hololive",
			Description: "A very special place for everyone",
			Logo: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//hoshimachi-suisei-hololive.gif",
			CategoryID: uuid.MustParse("928c8cc1-a75c-4fd2-a080-f2e4e090e698"),
			Banner: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(125).png",
			Rule: "Please support your oshis wholeheartedly",
			CreatorID: uuid.MustParse("ae2cf219-a6e3-46d7-aee0-6f05eb23f602"),
		},
		{
			Name: "Hoyoverse",
			Description: "A place that talks about the best game creator",
			Logo: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(2).png",
			CategoryID: uuid.MustParse("b49f4771-6207-46a2-84ef-8208775ae922"),
			Banner: "https://nyijguoyjemxmdyzzdry.supabase.co/storage/v1/object/public/registration-uploads//Screenshot%20(3).png",
			Rule: "Enjoy hoyoverse games",
			CreatorID: uuid.MustParse("f3122b31-7b92-400d-b94d-28ef100ca23b"),
		},
	}
	for _, element := range communities {
		db.Create(&element)
	}
	
	communities = generateRandomCommunities()
	for _, element := range communities {
		db.Create(&element)
	}


	// var hashtags = []model.Hashtag{
	// 	{
	// 		Name: "HSR",
	// 	},
	// 	{
	// 		Name: "Genshin",
	// 	},
	// 	{
	// 		Name: "Hololive",
	// 	},
	// 	{
	// 		Name: "Fate",
	// 	},
	// 	{
	// 		Name: "Graduation",
	// 	},
	// }
	// for _, element := range hashtags {
	// 	db.Create(&element)
	// }

	// var threadHashtags = []model.ThreadHashtag{
	// 	{
	// 		ThreadID: uuid.MustParse("df723ffb-8fb6-4e45-b43f-8a20f0cdce50"),
	// 		HashtagID: uuid.MustParse("510aed57-1728-4d89-8336-9c9f2e024bdf"), // HSR
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("df723ffb-8fb6-4e45-b43f-8a20f0cdce50"),
	// 		HashtagID: uuid.MustParse("c3113a07-a0c5-4fce-923b-bb261ef18b40"), // Genshin
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("df723ffb-8fb6-4e45-b43f-8a20f0cdce50"),
	// 		HashtagID: uuid.MustParse("76ef3273-7b93-4711-b9f6-af7ec49d4dfb"), // Hololive
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("e0fb98bf-6e5a-48fb-a46b-ba164d08af16"),
	// 		HashtagID: uuid.MustParse("598977df-daeb-404c-9874-8ad690d441c4"), // Fate
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("e0fb98bf-6e5a-48fb-a46b-ba164d08af16"),
	// 		HashtagID: uuid.MustParse("3fa46c5d-4e1b-4bd7-be29-46e6577eb889"), // Graduation
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("e0fb98bf-6e5a-48fb-a46b-ba164d08af16"),
	// 		HashtagID: uuid.MustParse("510aed57-1728-4d89-8336-9c9f2e024bdf"), // HSR
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("0cde263a-159c-4029-8006-dde96baa57b9"),
	// 		HashtagID: uuid.MustParse("c3113a07-a0c5-4fce-923b-bb261ef18b40"), // Genshin
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("0cde263a-159c-4029-8006-dde96baa57b9"),
	// 		HashtagID: uuid.MustParse("76ef3273-7b93-4711-b9f6-af7ec49d4dfb"), // Hololive
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("6b56d786-ed6b-4e40-851c-941424fa1860"),
	// 		HashtagID: uuid.MustParse("598977df-daeb-404c-9874-8ad690d441c4"), // Fate
	// 	},
	// 	{
	// 		ThreadID: uuid.MustParse("6b56d786-ed6b-4e40-851c-941424fa1860"),
	// 		HashtagID: uuid.MustParse("3fa46c5d-4e1b-4bd7-be29-46e6577eb889"), // Graduation
	// 	},
	// }
	// for _, element := range threadHashtags {
	// 	db.Create(&element)
	// }
}
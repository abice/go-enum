package main

import (
	"flag"
	"time"

	"github.com/mkideal/log"
	"gopkg.in/redis.v5"

	"github.com/mkideal/pkg/storage"
	"github.com/mkideal/pkg/storage/example/demo"
	"github.com/mkideal/pkg/storage/goredisproxy"
)

var (
	flRedisHost = flag.String("redis_host", "127.0.0.1:6379", "redis host")
	flRedisPwd  = flag.String("redis_pwd", "", "redis password")
)

func main() {
	flag.Parse()
	defer log.Uninit(log.InitColoredConsole(log.LvFATAL))
	log.SetLevel(log.LvDEBUG)

	client := redis.NewClient(&redis.Options{
		Addr:         *flRedisHost,
		Password:     *flRedisPwd,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	client.FlushDb()

	cache := goredisproxy.New(client)

	// NewEngine with name,database proxy and cache proxy
	eng := storage.NewEngine("test", storage.NullDatabaseProxy, cache)
	demo.Init(eng)

	// SetErrorHandler
	eng.SetErrorHandler(func(action string, err error) error {
		log.Printf(storage.ErrorHandlerDepth, log.LvWARN, "<%s>: %v", action, err)
		return err
	})

	api := eng.NewSession()
	defer api.Close()

	// Insert inserts records
	inserted := &demo.User{Id: 1, Name: "test1", Age: 10}
	api.Insert(inserted)

	// Get
	loaded := &demo.User{Id: 1}
	found, err := api.Get(loaded)
	log.Debug("Get: found=%v, error=%v, value=%v", found, err, loaded)

	// Update
	inserted.Age = 20
	api.Update(inserted, inserted.Meta().F_age)
	//loaded = &demo.User{Id: 1}
	//found, err = api.Get(loaded)
	//log.Debug("Get: found=%v, error=%v, value=%v", found, err, loaded)

	// Remove
	api.Remove(inserted) // <=> api.RemoveRecords(inserted.TableMeta(), inserted.Key())
	found, err = api.Get(&demo.User{Id: inserted.Id})
	log.Debug("Remove: found=%v, error=%v", found, err)

	users := []demo.User{
		{Id: 1, Name: "test1", Age: 10, AddrId: 1000, ProductId: 100},
		{Id: 2, Name: "test2", Age: 20, AddrId: 2000, ProductId: 200},
		{Id: 3, Name: "test3", Age: 30, AddrId: 3000, ProductId: 300},
	}
	products := []demo.Product{
		{Id: 100, Price: 1, Name: "p1", Image: "img1", Desc: "desc1"},
		{Id: 200, Price: 2, Name: "p2", Image: "img2", Desc: "desc2"},
		{Id: 300, Price: 3, Name: "p3", Image: "img3", Desc: "desc3"},
	}
	addresses := []demo.Address{
		{Id: 1000, Addr: "Beijing"},
		{Id: 2000, Addr: "Shanghai"},
		{Id: 3000, Addr: "Taiwan"},
	}
	keys := make([]int64, 0, len(users))
	for i := range users {
		api.Insert(&users[i])
		keys = append(keys, users[i].Id)
	}
	for i := range products {
		api.Insert(&products[i])
	}
	for i := range addresses {
		api.Insert(&addresses[i])
	}

	// Find
	us := demo.NewUserSlice(len(users))
	api.Find(storage.Int64Keys(keys), us)
	log.Debug("Find users: %v", us.Slice())

	// test ErrorHandler
	loaded = &demo.User{Id: 1}
	err = api.Update(loaded, "invalid_field")

	// FindView
	userview := demo.NewUserViewSlice(len(users))
	api.FindView(demo.UserViewVar, storage.Int64Keys(keys), userview)
	log.WithJSON(userview.Slice()).Debug("FindView: userview")

	// IndexRank
	rank, _ := api.IndexRank(demo.UserAgeIndexVar, 1)
	log.Debug("IndexRank: user %d UserAgeIndex: %d", 1, rank)
	rank, _ = api.IndexRank(demo.UserAgeIndexVar, -1)
	log.Debug("IndexRank: user %d UserAgeIndex equal to storage.InvalidRank? %v", -1, rank == storage.InvalidRank)

	// IndexScore
	score, _ := api.IndexScore(demo.UserAgeIndexVar, 1)
	log.Debug("IndexScore: user %d UserAgeIndex: %d", 1, score)
	score, _ = api.IndexScore(demo.UserAgeIndexVar, -1)
	log.Debug("IndexScore: user %d UserAgeIndex equal to storage.InvalidScore? %v", -1, score == storage.InvalidScore)

	// Clear
	api.Clear(demo.User{}.Meta().Name())
	api.Clear(demo.Address{}.Meta().Name())
	api.Clear(demo.Product{}.Meta().Name())
}

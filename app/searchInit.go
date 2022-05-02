package app

// func initializeAndRebuildIndex(config ezsearch.Config) {
// 	count, err := ezsearch.ProductIndex.DocCount()
// 	if err != nil {
// 		logger.Error("Failed while check intIndex count", err)
// 	}
// 	logger.Info(fmt.Sprintf("product index doc count is %d", count))
// 	if count == 0 {
// 		rebuildProductIndex(config, catalogboltdb.JsonDb, ezsearch.ProductIndex)
// 	}
// 	rebuildIntIndex(config)

// 	count, err = ezsearch.AggProductIndex.DocCount()
// 	if err != nil {
// 		logger.Error("Failed while check agg product Index count", err)
// 	}
// 	logger.Info(fmt.Sprintf("agg product index doc count is %d", count))
// 	if count == 0 {
// 		rebuildAggProductIndex(config, catalogboltdb.JsonDb, ezsearch.AggProductIndex)
// 	}
// }

// func rebuildIntIndex(config ezsearch.Config) rest_errors.RestErr {
// 	recCount, err := ezsearch.IntIndex.DocCount()
// 	if err != nil {
// 		logger.Error("Failed while check intIndex count", err)
// 		return rest_errors.NewInternalServerError("Failed while check intIndex count", err)
// 	}
// 	logger.Info(fmt.Sprintf("intergration index doc count %d", recCount))
// 	if recCount == 0 {
// 		logger.Info("Rebuild intergration index ...")
// 		intData, err := services.CreateOrUpdateProductIntService.GetAllProductexternalInt() //TODO rebuild logic to load data page avoid performacne issues
// 		if err != nil {
// 			logger.Error("Failed while fetch CreateOrUpdateProductIntService.Search", err)
// 			return rest_errors.NewInternalServerError("Failed while fetch CreateOrUpdateProductIntService.Search", err)
// 		}

// 		i := ezsearch.IntIndex
// 		count := 0
// 		startTime := time.Now()
// 		logger.Info("Indexing integration...")
// 		batch := i.NewBatch()
// 		batchCount := 0
// 		for _, item := range intData {
// 			m := ezsearch.SearchIntModelClient{ProductSku: item.ProductSku, LanguageCode: item.LanguageCode,
// 				ExternalIntId: item.ExternalIntID, ExternlRefId: item.ExternlRefID}
// 			var err error
// 			mj := models.GetMasterScheduleJobByLangCode(item.LanguageCode)
// 			//fmt.Println("mj := models.GetMasterScheduleJobByLangCode(item.LanguageCode)", item.LanguageCode)
// 			key := fmt.Sprintf("%s_%s_%s", mj.CatalogType, item.ProductSku, item.LanguageCode)
// 			catBytes, err := catalogboltdb.GetKey("" , key)
// 			if err != nil {
// 				logger.Error(fmt.Sprintf("Failed while access catalog|rebuildIntIndex %s", key), err)
// 			} else {
// 				var catM models.OccProdutClient
// 				err = json.Unmarshal(catBytes, &catM)
// 				if err != nil {
// 					logger.Error(fmt.Sprintf("Failed while unmarshal|rebuildIntIndex %s", key), err)
// 				} else {
// 					m.Category = strings.Join(catM.ParentCategoryPath, "|")
// 					m.Classification = catM.DefinitionName
// 				}
// 			}

// 			//("item.creteddate", item.CreatedAt, item.ModifiedDate)
// 			m.CreatedDate = item.CreatedDate
// 			// if err != nil {
// 			// 	logger.Error("Failed whilte parse time|CreatedDate", err)
// 			// }
// 			m.ModifiedDate = item.ModifiedDate
// 			// if err != nil {
// 			// 	logger.Error("Failed whilte parse time|ModifiedDate", err)
// 			// }

// 			//fmt.Println("dateddddd", m.CreatedDate, m.ModifiedDate, item.CreatedDate, item.ModifiedDate)
// 			key1 := fmt.Sprintf("%s_%s_%s", item.ProductSku, item.ExternalIntID, item.LanguageCode)
// 			batch.Index(key1, m)
// 			batchCount++

// 			if batchCount >= config.IndexBatchSize {
// 				err := i.Batch(batch)

// 				if err != nil {
// 					logger.Error("Failed while update|i.Batch", err)
// 					continue
// 				}
// 				batch = i.NewBatch()
// 				batchCount = 0
// 			}
// 			count++
// 			if count%1000 == 0 {
// 				indexDuration := time.Since(startTime)
// 				indexDurationSeconds := float64(indexDuration) / float64(time.Second)
// 				timePerDoc := float64(indexDuration) / float64(count)
// 				log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
// 			}
// 		}
// 		// flush the last batch
// 		if batchCount > 0 {
// 			err := i.Batch(batch)
// 			if err != nil {
// 				logger.Error("Failed while last batch", err)
// 			}
// 		}
// 	}
// 	return nil
// }

// func rebuildProductIndex(config ezsearch.Config, db *bolt.DB, i bleve.Index) {
// 	countMap := make(map[string]int)
// 	db.View(func(tx *bolt.Tx) error {
// 		count := 0
// 		startTime := time.Now()
// 		logger.Info("Indexing product...")
// 		batch := i.NewBatch()
// 		batchCount := 0

// 		c := tx.Bucket([]byte("DB")).Bucket([]byte(config.BoltDbBucketName)).Cursor()
// 		for k, v := c.First(); k != nil; k, v = c.Next() {
// 			kPart := strings.Split(string(k), "_")
// 			//fmt.Println("Lang Fileter", kPart[len(kPart)-1], hybris.LangFilters)

// 			lang := kPart[len(kPart)-1]
// 			if countMap[lang] == 0 {
// 				countMap[lang] = 0
// 			}
// 			countMap[lang] = countMap[lang] + 1
// 			//fmt.Println("lang code", lang, countMap[lang])
// 			if !hybris.LangFilters[lang] {
// 				continue
// 			}
// 			var m models.OccProdutClient
// 			//fmt.Printf("indexing key=%s\n", k)
// 			err := json.Unmarshal(v, &m)
// 			if err != nil {
// 				logger.Error("Failed unmarshal|updateIndex", err, zap.String("key", string(k)))
// 			}
// 			//m.ID = strings.Replace(m.ID, "-", "_", -1)

// 			dispName := fmt.Sprintf("%v", m.DisplayName)

// 			if len(dispName) > 0 {
// 				val := fmt.Sprintf("%v", m.DisplayName)
// 				valpart := strings.Split(val, ":")
// 				if len(valpart) > 1 {
// 					dispName = strings.Replace(valpart[1], "]", "", -1)
// 				} else {
// 					dispName = valpart[0]
// 				}

// 			}

// 			if len(m.PropertyBag[global.DISCONTINUEDDATE]) > 0 {
// 				m.PropertyBag[global.DISCONTINUEDDATE] = m.PropertyBag[global.DISCONTINUEDDATE][:19]
// 			}
// 			if len(m.PropertyBag[global.MASTERLAUNCHEDDATE]) > 0 {
// 				m.PropertyBag[global.MASTERLAUNCHEDDATE] = m.PropertyBag[global.MASTERLAUNCHEDDATE][:19]
// 			}
// 			if err != nil {
// 				logger.Error("Failed while conver launched date", err)
// 			}

// 			batch.Index(string(k), m)

// 			batchCount++

// 			if batchCount >= config.IndexBatchSize {
// 				err = i.Batch(batch)

// 				if err != nil {
// 					logger.Error("Failed while update|i.Batch", err, zap.String("key", string(k)))
// 					continue
// 				}
// 				batch = i.NewBatch()
// 				batchCount = 0
// 			}
// 			count++
// 			if count%1000 == 0 {
// 				indexDuration := time.Since(startTime)
// 				indexDurationSeconds := float64(indexDuration) / float64(time.Second)
// 				timePerDoc := float64(indexDuration) / float64(count)
// 				log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
// 			}
// 		}

// 		// fmt.Println("countMap:", countMap["en-US"])
// 		// fmt.Println("countMap:", countMap["en-IE"])
// 		// fmt.Println("countMap:", countMap["en-GB"])
// 		// flush the last batch
// 		if batchCount > 0 {
// 			err := i.Batch(batch)
// 			if err != nil {
// 				logger.Error("Failed while last batch", err)
// 			}
// 		}

// 		return nil
// 	})
// 	//for k, v := range countMap {

// 	//}
// }

// func rebuildAggProductIndex(config ezsearch.Config, db *bolt.DB, i bleve.Index) {
// 	countMap := make(map[string]int)
// 	db.View(func(tx *bolt.Tx) error {
// 		count := 0
// 		startTime := time.Now()
// 		logger.Info("Indexing aggragated product...")
// 		batch := i.NewBatch()
// 		batchCount := 0

// 		c := tx.Bucket([]byte("DB")).Bucket([]byte(config.BoltDbBucketName)).Cursor()
// 		for k, v := c.First(); k != nil; k, v = c.Next() {
// 			kPart := strings.Split(string(k), "_")
// 			//fmt.Println("Lang Fileter", kPart[len(kPart)-1], hybris.LangFilters)

// 			lang := kPart[len(kPart)-1]
// 			if countMap[lang] == 0 {
// 				countMap[lang] = 0
// 			}
// 			countMap[lang] = countMap[lang] + 1
// 			//fmt.Println("lang code", lang, countMap[lang])
// 			if !hybris.LangFilters[lang] {
// 				continue
// 			}
// 			var m models.OccProdutClient
// 			//fmt.Printf("indexing key=%s\n", k)
// 			err := json.Unmarshal(v, &m)
// 			if err != nil {
// 				logger.Error("Failed unmarshal|updateIndex", err, zap.String("key", string(k)))
// 			}
// 			//m.ID = strings.Replace(m.ID, "-", "_", -1)

// 			dispName := fmt.Sprintf("%v", m.DisplayName)

// 			if len(dispName) > 0 {
// 				val := fmt.Sprintf("%v", m.DisplayName)
// 				valpart := strings.Split(val, ":")
// 				if len(valpart) > 1 {
// 					dispName = strings.Replace(valpart[1], "]", "", -1)
// 				} else {
// 					dispName = valpart[0]
// 				}

// 			}
// 			sm := ezsearch.AggregateProductClient{}
// 			sm.DisplayName = dispName
// 			sm.CatalogId = m.CatalogID
// 			sm.LanguageCode = m.LangCode
// 			sm.Sku = m.ID
// 			sm.ContentModifiedDate = m.ContentModifiedDate
// 			if len(m.ContentModifiedDate) == 0 {
// 				sm.ContentModifiedDate = time.Now().UTC().Format(date_utils.UTCDateLayout)
// 			}

// 			sm.Status = m.PropertyBag[global.PRODUCTSTATUS]
// 			sm.Classification = m.DefinitionName
// 			// if m.ID == "DCD710D2-KS" {
// 			// 	fmt.Println("class", sm.Sku, sm.LanguageCode, sm.Classification, m.DefinitionName)
// 			// 	fmt.Println("prd str", string(v))
// 			// 	fmt.Println("prd model", m)

// 			// }

// 			sm.Category = strings.Join(m.ParentCategoryPath, " ")
// 			sm.Warrantycode = m.PropertyBag[global.WARRANTYCODE]
// 			sm.Key = string(k)

// 			if len(m.PropertyBag[global.DISCONTINUEDDATE]) > 0 {
// 				sm.Discontinueddate = m.PropertyBag[global.DISCONTINUEDDATE][:19]
// 			}
// 			if len(m.PropertyBag[global.MASTERLAUNCHEDDATE]) > 0 {
// 				sm.Launcheddate = m.PropertyBag[global.MASTERLAUNCHEDDATE][:19]
// 			}
// 			ezsearch.SetAggProductIndexIntDates(&sm, m.ID, m.LangCode)

// 			batch.Index(string(k), sm)

// 			batchCount++

// 			if batchCount >= config.IndexBatchSize {
// 				err = i.Batch(batch)

// 				if err != nil {
// 					logger.Error("Failed while update|i.Batch", err, zap.String("key", string(k)))
// 					continue
// 				}
// 				batch = i.NewBatch()
// 				batchCount = 0
// 			}
// 			count++
// 			if count%1000 == 0 {
// 				indexDuration := time.Since(startTime)
// 				indexDurationSeconds := float64(indexDuration) / float64(time.Second)
// 				timePerDoc := float64(indexDuration) / float64(count)
// 				log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
// 			}
// 		}

// 		// fmt.Println("countMap:", countMap["en-US"])
// 		// fmt.Println("countMap:", countMap["en-IE"])
// 		// fmt.Println("countMap:", countMap["en-GB"])
// 		// flush the last batch
// 		if batchCount > 0 {
// 			err := i.Batch(batch)
// 			if err != nil {
// 				logger.Error("Failed while last batch", err)
// 			}
// 		}

// 		return nil
// 	})
// 	//for k, v := range countMap {

// 	//}
// }

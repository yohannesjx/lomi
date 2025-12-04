You are an expert Go backend engineer building a TikTok-style short-video + live-streaming app called "Lomi social".

I have:
- My current production Go backend in the folder "lomi_mini" dating app with auth, profiles, coins, payments, real-time chat, PostgreSQL...
- Full source code of native Android (Kotlin) + iOS (Swift) TikTok clone apps bought from CodeCanyon
- Complete reverse-engineered API contract in the attached file "TIKTOK_API_CONTRACT.md"

I am discarding the PHP backend completely. I want to use ONLY my existing Go backend.

Task:
1. Analyze my current "lomi_mini" Go code and database structure.
2. Implement the exact 6 critical endpoints from the contract so the Android/iOS apps work perfectly:
   • POST /api/registerUser
   • POST /api/showUserDetail
   • POST /api/showRelatedVideos
   • POST /api/liveStream
   • POST /api/sendGift
   • POST /api/purchaseCoin

3. For live streaming: use mediamtx or any better

4. Reuse as much of my existing code as possible:
   - Auth & JWT middleware
   - User model & profile logic
   - Coin/wallet system
   - Payment flow

5. Strictly match the exact JSON request/response format from TIKTOK_API_CONTRACT.md (including "code":200, "msg": {...})

6. Output:
   - One complete file: routes_streaming.go (ready to drop into my project)

7. Bonus (if easy): add a simple dummy /api/showRelatedVideos that returns 5 fake videos so the feed loads immediately.

I only have 1 server right now and will add more later — keep everything simple and scalable.

Start coding now. Use only standard Go libraries whatever I’m already using and you can add whatever is best practice. No external packages unless absolutely necessary.


---
title: "Bitcoin and economic nihilism"
date: 2022-10-10
series: ethereum
tags:
 - cryptocurrency
 - web3
 - blockchain
 - bitcoin
 - nft
---

<xeblog-hero ai="Waifu Diffusion v1.2" file="inkopolis-sunset" prompt="anime style, matte painting, splatoon, giant sword, squirt guns, cityscape, skyline, sunset, concert, studio ghibli, serene, peaceful"></xeblog-hero>

I have picked a hell of a time to experiment with Ethereum. In general, I've been trying to figure out for myself if the whole blockchain scene is something that is actually worth the scorn and toxicity I've had for it for so long. This technology is _really neat_ from a technical standpoint, but being a cool nerd toy and a viable product for the masses are polar opposites. Especially when money is involved.

I want to start out by saying that I am not a neutral party here. I have been anti-blockchain [to the point of toxicity](https://xeiaso.net/blog/trying-to-use-security-token) for a long time. I have also been harmed by the [Mt. Gox exit scam](https://moneybadger.stocktwits.com/history-of-mt-gox-bitcoin/). When that happened, I was in college. My father was paying a lot of my college tuition by selling the Bitcoin he had mined back when the project was brand new and CPU mining was viable (I think he mined like 4 blocks or something). He used Mt. Gox to exchange that Bitcoin for USD to pay the school. The exit scam happened while I was at college. Thankfully it didn't happen just before one of the quarter due dates, but I'm pretty sure it was a large part of the reason why I ended up having to drop out of college. That and unmedicated ADHD, but that's a story for another day.

In the past, I have been toxic to people doing blockchain stuff out of that past trauma re-emerging. I've needlessly cut off friendships. I don't want to be known for toxicity anymore. As part of this healing process, I've been trying to _actually understand_ this technology so I can figure out if this technology is worth using or not. I don't want to rely on other people's understanding anymore. If I do eventually conclude that there is something there, I'm willing to admit that I was wrong. This is part two of my magical blockchain journey. You can read the first part [here](https://xeiaso.net/blog/trying-to-use-security-token) or see the entire series [here](https://xeiaso.net/blog/series/ethereum).

<xeblog-conv name="Open_Skies" mood="wave">Hi! I’m Open, the technical editor for this blog. Usually I’m in the background rather than the article itself. I’m working with Xe on a Nix flakes post which was supposed to be the first one I formally co-author, but I had some Thoughts about this topic so I’m jumping in here.</xeblog-conv>

Today I want to investigate the ways Ethereum was sold to _me_, my thoughts on them as I've been researching things, and some of the other philosophical points of cryptocurrency in general that are worth consideration. As another one of the biases I have going into this, I would love to use something like Ethereum to be able to pay for commissioning art without having to worry about payment processors like PayPal canceling accounts for that art being erotica. Ideally, I'd like to be able to give the artist cash, but that doesn't work so well when the only contact I have with them is over the Internet.

It would be nice if there were other uses too, but really the big selling point in my book is the ability to make payments to people without a payment processor being able to say "no" and ban you because you want something perfectly legal that makes puritans squirm (such as erotica).

I would love for my past assumptions about the technology to be proven wrong. If everything pans out like they think it will, that could be great. Keywords _could_ and _if_.

<xeblog-sticker name="Cadey" mood="coffee"></xeblog-sticker>

## The promises I was told

<xeblog-conv name="Cadey" mood="coffee">Subjectivity warning: remember that I am talking about the things _I personally_ was promised as assets to the cryptocurrency space and/or reasons why I should care about this at all. This is subjective. Your experiences can and will be different.</xeblog-conv>

When the web3 people have tried to sell me on the web3 stuff they've been working on, they usually focus on these high level things:

- Peer-to-peer payments without a payment processor or a bank being able to say "lol no"
- Anonymous payments, so that I have privacy the same way that I'd have if using cash
- Unfakeable "ownership" of arbitrary tokens so that I can prove that I have things like furry OCs or the "right to use a closed species"
- Distributed computing, so that I can make additional bits of code run to validate things involving my presence on the blockchain, create arbitrary tokens (fungible and non-fungible) and more
- Being able to use your hat or avatar in another game through blockchain cross-pollination of assets

To be honest, this is a huge thing, and if even the first one is viable then it
means that I can have less fear about making transactions on the Internet. I
could use this to commission art, help friends meet rent when they land on hard
times, and not have to worry about someone changing their name on all of their
paperwork everywhere all at once so that we can send money to each other. I
could just send them money to their Ethereum addresses. They could just send
money to `xeiaso.eth` or `0xeA223Ca8968Ca59e0Bc79Ba331c2F6f636A3fB82`.
Everything will work as normal.

<div class="warning"><xeblog-conv name="Cadey" mood="coffee">Due to how this all
is designed, anyone can send anything into anyone's cryptocurrency wallets.
Beware anyone trying to call me out for things that have been shoved into my
wallet against my will.</xeblog-conv></div>

<xeblog-conv name="Numa" mood="delet">Except when it comes to paying rent, then you will need a middleman to convert your internet dollars into rent paying dollars.</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Yeah. This is where the golden promised future of cryptocurrency starts falling apart. It's like US dollars in a way. You can only use US dollars with people that accept US dollars as payment. You can only use Ethereum with people that accept Ethereum as payment. There's a huge network effect at play here, and in order to have any cryptocurrency become a viable means of payment, you have to play the massive uphill game of trying to establish a network effect of your own. Ideally you want taking Ethereum to be _easier_ than taking US dollars. Guess what isn't the case.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="wave">For me, the whole dance of paying someone for content on the internet is… rough. Like, there are so many steps to pay someone, and everyone wants to send me emails or sign me up for a subscription to pay every month for premium content, and I need to scroll around the page to find a payment link then go through multiple forms and find my credit card to enter the number or sign into paypal or whatever. I’d really like to be able to just like… hey, this research you published for free on your blog is excellent and saved me like three hours of trying to put it together from other sources. One click, you get a dollar, thank you very much, I’ll be on my way. And on the other hand, I’d like to be able to have something like that which is trivial to set up on my own content rather than needing advertisement peering and accounts and having a minimum payout before I can use the money, and annoying all the viewers…</xeblog-conv>

<xeblog-conv name="Cadey" mood="enby">Amazingly enough, a lot of this is actually solved in Canada provided the person or business you're trying to send money to has a bank account. You can use [Interac](https://en.wikipedia.org/wiki/Interac) to send money to people by phone number or email address. It's how I've paid my tax attorney and for a while it's how I've paid rent. This is only viable for bigger payments though, something like [micropayments](https://en.wikipedia.org/wiki/Micropayment) have never really taken off because of transaction fees, especially with cryptocurrencies. Litecoin has [very low transaction fees](https://litecoin.info/index.php/Transaction_fees), but if your total payment is a dollar then a 4% transaction fee sounds kind of ludicrous. If things like [Stellar](https://en.wikipedia.org/wiki/Stellar_(payment_network)) were more viable, then it would be a different story.<br /><br />Until then we'll be dealing with Twitch only paying out once you make $100 or my advertising provider only paying out when you make $50. It's a nightmare.</xeblog-conv>

<xeblog-conv name="Numa" mood="stare">Wait, wait. One of the original points in the bitcoin paper was that it would avoid paying transaction fees to the bigger banks. When you pay a transaction fee on a cryptocurrency payment, where does it go? The miners/stakers? Wouldn't that just end up making you pay the transaction fees to the bigger players anyways?</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Yep. Stay tuned for why.</xeblog-conv>

In practice, these promises don't pan out as easily as they sound. You actually do need middlemen. A lot of them, especially if you want to pay rent or buy food with that money. You'll need to register with one of many "exchanges" which may or may not have the best security practices while at the same time demanding copies of your ID documents so they can comply with "know your customer" laws. That seems safe.

You can send payments to artists to commission art with Ethereum, but the main problem there is that transactions on the blockchain aren't reversible in cases of fraud. It seems that the Ethereum people are of the standpoint that the technology of the blockchain and the cryptosystem make fraud categorically impossible in Ethereum land. This is...a take, but in general I can kind of see what they are going for. Reversibility of transactions is usually a byproduct of the traditional financial system having middlemen that can make calls like "yeah this person was out of line" and then undo the transaction (and even restrict the account of the perpetrator) so the victim has an easier time and they can reduce the amount of fraud on their network.

<xeblog-conv name="Open_Skies" mood="snug">It’s not just a middleman making a call that the person was out of line, it’s that the traditional financial system has the traditional legal system with the ability to go after someone who committed fraud or theft and force them to pay the money back afterwards (at least in theory). That ability is what enables payment processors to offer payment protection features, without being escrow agencies. And it also means that the middlemen aren’t required for that protection, because individuals can sue for payment fraud or theft. The US Bill of Rights guarantees the right to a lawsuit with trial by jury when more than $20 is in question, and anyone can file a suit; the middlemen just have the ability to threaten to stop providing payments for all customers if the disputed payment isn’t returned, the staff and premade forms that make it very easy for them to file such suits when necessary, and enough money to absorb chance without difficulty. And it’s not like reversing a transaction is *impossible* in Ethereum, but it’s a lot harder, to the point that the very few times it has ever happened were ecosystem-altering events.</xeblog-conv>


<xeblog-conv name="Mara" mood="hmm">So this ecosystem where payments are irreversible by design has sometimes had the community come together to reverse payments they thought were fraudulent? Seems legit...</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">Yep! And we can talk about the story there later! It’s interesting. Anyways, there is some stuff to deal with fraud on Ethereum without the… events, like Kleros, which provides a decentralized anonymous trial-by-jury approach to dealing with disputes. It could, theoretically, be integrated with a kind of escrow system where both parties pay into a contract that releases their payment/item after a set time or depending on the judgment of a dispute, but if anyone has done this, it is certainly not widespread enough for most transactions to have that kind of security. Also, it adds complexity to the transaction, meaning more gas is required to perform it, which means the transaction costs more. Whether or not it costs more than a credit card’s overhead we’ll have to look at later, because we need more concepts before we can go through the math to see how much transactions cost.</xeblog-conv>

In Ethereum land, identity is determined by being able to sign messages with a private key. Ownership of that private key gives you ownership of the wallet, and by extension everything _in_ that wallet. The entire system is based on the assumption that humans can keep track of private keys and keep them safe. Without backing them up on things like iCloud, Google Drive, Amazon S3, or any other online storage system. People usually suggest you put a backup copy of your recovery phrases (private keys) on paper (written in pencil, because apparently pen ink fades over longer periods of time) inside a safety deposit box. This is surely easy for people that are more affluent, but for ordinary people that is an extra $40 per month fee, which can be more than twice what an account fee is for normal checking accounts.

<xeblog-conv name="Cadey" mood="coffee">It would be fine if these processes encouraged the use of splitting and backing up keys in ways that align with best practices across the industry, but they don't. They suggest putting the entire keys to your digital kingdom into a single safety deposit box. There are devices that allow you to do hardware escrow of cryptocurrency private keys, but again they still encourage you to back up the key on paper. They also act like safety deposit boxes are free. They are not. They are expensive.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">Ok, but “best practices across the industry” are kinda shitty in a lot of places.</xeblog-conv>

Ethereum addresses are anonymous.


<xeblog-conv name="Open_Skies" mood="idle">More like pseudonymous.</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Okay, granted, but Ethereum addresses are 160 bit values that are not fundamentally a human's name.</xeblog-conv>

However, transactions you make with those addresses are _not_ hidden. You can look at the [entire transaction history of any account you want](https://etherscan.io/address/0xea223ca8968ca59e0bc79ba331c2f6f636a3fb82), which makes choosing targets for phishing attacks trivial. If you can associate a person with an address, you can know what they use their money with. You can figure out how much of a target they are to you. If your identity gets tied to the wrong address, that address is blown and you need to potentially abandon that wallet and start over with a new one.

And then, of course there are the one-two punch of smart contracts and non-fungible tokens. I don't feel confident talking about smart contracts yet, so I'm going to defer talking about them until a later article. But let's take a look at non-fungible tokens a bit closer. Those are fun.

<xeblog-conv name="Mara" mood="hmm">Fun…how?</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="snug">I can talk about smart contracts though, and you do kinda need to know what a smart contract is to know what an NFT is, because smart contracts are what give NFTs meaning, to the extent they have meaning at all. After all, it's all just bits, and bits don't have color.</xeblog-conv>

## Non-fungible tokens, CS:GO knife skins and furry art

<xeblog-hero ai="Stable Diffusion v1.4" file="teal-yellow-utility-poles" prompt="teal and yellow colors. Utility poles in style of cytus and deemo, mysterious vibes, set in half-life 2, beautiful with eerie vibes, very inspirational, very stylish, surrealistic, perfect digital art, mystical journey in strange world, bastion game"></xeblog-conv>

[Non-fungible tokens](https://ethereum.org/en/nft/) are little bits of data that signify that some wallet has some bit of data. This data is not able to be forged and is independently verifiable on the Ethereum blockchain. There is insufficient court precedent to know what having an NFT associated with your cryptocurrency address actually means, but to give people the benefit of the doubt I'm going to assume that you own the _token itself_, not any of the data associated with it such as images.

<xeblog-conv name="Open_Skies" mood="snug">What does "owning the token" mean though? Ultimately all ownership must come down to some form of exclusion. In order to say that you own something you must have some way to say that someone else can't do something because they don't own it. Some NFT people try to say things like that owning the NFT provides an exclusive right to… have the image saved on your device? Which is spectacularly stupid. The asset linked to in a NFT lives on the original creators' hard drive, and people don't claim (as far as I've seen) that the creator must immediately delete their own copy of it once the NFT sells. Additionally, most NFT assets are hosted on IPFS, which specifically, deliberately, by its very operation, copies the asset to many different storage devices as people retrieve and pin it. Most browsers store temporary caches of every image they retrieve to speed up page loads. This criticism of NFTs actually made me anti-anti-crypto for a while. "Haha, we're going to right click and save your NFT" just seems like such a spectacularly stupid criticism of the system that it couldn't possibly connect to any real person that was willing to put down money on an immature technological R&D effort; surely the people making that claim are acting in bad faith, right? Right?</xeblog-conv>

<xeblog-conv name="Numa" mood="delet"><br />&gt;humans<br />&gt;rational<br /><br />pick one.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="snug">Turns out, nope. NFT people got really upset about people doing that, and tried to coin "right-clickers" as a derogatory term and threatened IP litigation for privately downloading publicly posted content deliberately made available for download. And seeing that happen, seeing people I thought couldn't exist, that talking about them had to be in bad faith, come out of the woodwork and turn out to be either the majority of the community or actively exploited and scammed by people pretending to believe that in bad faith, is what made me realize that NFTs are doomed to uselessness. Even if I can show that there are particular sets of useful and pro-social desiderata that can only be solved by something isomorphic to NFTs, such a system can't be built because they're incompatible with modern humans.</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">This really sucks to be honest. Something like this could give us a lot of freedom when it comes to designing and implementing things. Such a decentralized data store could make it a lot easier to import art around. Imagine if you could fire up your VR home world and the paintings on the virtual walls were automatically lit up with the images attached to NFTs you own. Something where you could easily move all that data around to different applications. This is a nice _vision_, but there is a huge difference between a _vision_ and a viable product. It would be really cool to have an e-ink screen that could be used to show off digital art, but I don't really see this happening unless that screen is sabotaged to only show off NFT art. I'm pretty sure that people don't want to see the horrific abominations that are [mutant apes](https://opensea.io/collection/mutant-ape-yacht-club) (the sad results of overconsumption of slurp juices). This kind of a vision of the future is the kind of thing you only see in Zuckerberg keynotes, not actual viable things that can exist in this capitalist hellscape we live in.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="snug">Anyways, moving on from the rant, the most basic thing that owning an NFT gives you is bragging rights. You own the NFT, and others don't, so you get to brag about owning the NFT and others don't. Anything beyond that is modular and separate. For example, something that NFTs can and have been used for is tickets to an event; proving that you have a ticket cryptographically means that you can walk into the event, and it might arguably be superior to current ticket mechanisms in some ways, like allowing secondary resales in a guaranteed way. With physical tickets, you need to carry around a collection of physical things that could get shuffled around or dropped. Current electronic tickets, you are emailed a QR code that you can back up and save however you please, making it resistant to getting lost. However, you can't reliably resell an electronic ticket, because the seller still has the QR code and might get in line before you, thus invalidating the sale after the fact. On the one hoof, this means that you can recover money even if the event has a bad refund policy; on the other hoof it means people can scalp tickets readily.</xeblog-conv>

### Steam inventory

This is similar to the Steam inventory system. This lets you have items from games like Team Fortress 2 (TF2, a once widely popular war-themed hat simulator) show up in a global inventory. This lets other people on Steam see what special items or hats you have so they can trade them with you.

<xeblog-conv name="Cadey" mood="enby">The "mannconomy" in TF2 enabled me to buy games on the cheap during Steam sales while I was a starving college student. I used to have a couple virtual machines set up to run TF2 in separate accounts while I was in class. These virtual machines would continuously idle in specially set up maps and collect random weapon drops with an autohotkey loop. These weapons were then traded to me to turn into refined metal, which I used to sell for lootbox keys, which I then used to sell for money on PayPal so I could buy games. I don't expect this to be viable anymore with the [collapse of the "mannconomy"](https://www.thegamer.com/team-fortress-2-economy-collapse-hats/) after valuable hats got duplicated and caused hyperinflation of the refined metal to key ratio. I got a bunch of games that way though.</xeblog-conv>

Realistically: TF2, [CS:GO](https://store.steampowered.com/app/730/CounterStrike_Global_Offensive/) and [DotA 2](https://www.dota2.com/home) are the only games that really take full advantage of this system. It never really caught on outside of games that Valve made.

<xeblog-conv name="Open_Skies" mood="wave">What about, to pick some random examples, Redout: Enhanced Edition, Space Engineers, Battle block Theater, Don't Starve Together…</xeblog-conv>


<xeblog-conv name="Cadey" mood="coffee">Okay, granted, but how many of these games were there? And how many of those games are really that popular? Trading cards don't count, I can't use my Steam trading cards in games.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">[Several hundred.](https://steamcommunity.com/sharedfiles/filedetails/?id=873140323)</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Ok, but there was never a cross-pollination of items in one game being usable in another, right? You couldn't take your custom knife skin in CS:GO into DotA 2 to use on Bloodseeker's blades. You can't take your "Strange" TF2 hats into CS:GO so your counter-terrorist squad can be dripped out and fly. Each item really stayed in their own world. It was great for trading items though. Steam support even used to be able to reverse transactions if you got scammed.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="wave">What about Garry's Mod? What about all the hats that are wearable in both TF2 and Portal 2, or the other cosmetics that you get in Portal 2 from other games?</xeblog-conv>

<xeblog-conv name="Cadey" mood="facepalm">Right, those. Honestly the fact that it's not a hugely memorable feature should be a sign that it's not a market demand. And we're not even talking about the difficulty of putting arbitrary assets into arbitrary engines yet.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">I’ll grant you that the Portal 2 crossover hats are forgettable (I had to look it up to make sure I wasn’t misremembering (and I was in fact misremembering some things)), but not Garry’s Mod (gmod). Garry’s Mod’s ability to use assets from any Source engine game the user has purchased on Steam is a huge selling point. It’s prominently advertised and widely used, though it does suffer somewhat from requiring the consumer to own a license to the games the assets come from rather than the games implicitly coming with a license to distribute the assets for use in gmod maps, but it’s understandable why they don’t; that would be a kinda thorny IP case to try to navigate, and it would reduce the incentives people would have to buy your game.</xeblog-conv>

<xeblog-conv name="Cadey" mood="enby">True, Garry's Mod is why I ended up buying Counter Strike Source. I have no interest in playing Counter Strike but I was tired of seeing error objects and black/magenta textures everywhere. Oh god that texture is an _eyesore_.</xeblog-conv>

It's also widely understood that having a hat, weapon skin, announcer voice, or whatnot in your Steam inventory didn't mean you _owned the intellectual property_ of that item. You just had a unique copy of it in your inventory, much like having a Pokemon card didn't mean you owned the concept of Pikachu.

<xeblog-conv name="Open_Skies" mood="snug">Oh, another thing. SteamVR Home has a lot of interesting ways to obtain assets. It has scavenger hunts for official assets that encourage exploring a variety of worlds to see what the system has to offer, ways to win them in minigames, workshop integration for sharing user generated content, etc. But the things that are relevant to the topic at hand are that it has collectibles you can find in other games and decorate your virtual room or your “avatar”, such as it is. For example, I can put on a mask from Payday 2 in SteamVR because once upon a time some people promised me it was a good party experience and babysat/coached me through a few high level missions. I can line up a figurine of atlas from portal 2 next to an enemy from super hot, and shoot at them with a little remote controlled drone I won in a shooting range minigame. There are only about 70-80 games on Steam that actually provide SteamVR Collectibles, but that’s very different from zero. And of course, these items are just things in the Steam inventory so you can trade them. Also, you can bring Steam achievements from any game in as 3d trophies you can put in a virtual shelf or trophy case.</xeblog-conv>

### Furry "adoptables" and "closed species"

The furry fandom's social norms are a bit weird to outsiders. One of the big things that people have trouble with is the fact that there's no actual "source material" for the furry fandom. It's entirely self-standing and dates back to cartoonists wanting to draw funny animals. There's no real cohesive universe for the furry fandom, but we all like to pretend there is for our own collective amusement.

One of the things you will see people do in the furry fandom is create their own Original Characters (OCs). The rationale for these can boil down to a number of things, including roleplay, characters to use in fiction, sexual fantasies, an idealized version of the self, or any number of other things you can do with characters. For example, here's a shark typing on a split keyboard that you've seen before:

<xeblog-sticker name="Mara" mood="hacker"></xeblog-sticker>

I use these characters to help me illustrate side ideas, internal disagreements, and more on this blog. This cast of characters helps me use common tropes to portray ideas like "the student that always wants to learn more", "the author self-insert", or "the ruthless shitposter that really does mean well". I can also use them to have guests peek into the article and give their own input.

Most people end up creating their own characters out of their own hopes, ambitions, and dreams (such as the black-tipped reef shark on the blog with the cool keyboard that all those tech bros are envious of). Some artists have a lot of free time and decide to create additional characters. These artists will then auction them off on social media (yes, really) and hand over the reference sheet (an image that describes how the character _should_ look like from several angles and the color palette at play) once the winner pays.

From then on, the "owner" of the character is widely inferred to be the winner of the auction. This isn't actually verified anywhere in most cases, but using someone else's character without permission is a very harsh taboo to break. If people really want to verify, they have to track down the original artist and ask them. This is a decentralized system, but it usually works out because most people are actually honest about the "ownership" of adoptable characters.

Sometimes artists have ideas for an entire species of characters with its own lore, art and norms. A lot of the time, these ideas will be posted publicly for anyone to share and use. This is how you get things like [Sergals](https://en.wikifur.com/wiki/Sergal), which are sort of cheese-headed bipedal furry things. In the Sergal lore, they are native to `$SOME_OTHER_PLANET` and the creators place no limits on who can create and use Sergal characters as long as they follow the few core differentiators (they are not mammals, etc).

Other times artists aren't so open with their ideas. They may want to have an exclusive club for the people that "own" characters of that species. So, they "close" the species and say that only people they approve of can make characters with those characteristics. Then they charge a "licensing fee" for entry into that club. Breaking the taboo and making your own closed species character without a "license" can result in a lot of artists denylisting you and then it will make getting art of any of your characters difficult.

<xeblog-conv name="Cadey" mood="facepalm">Yes, this is a thing that actually happens in real life. I was almost subjected to this because the original idea for this character was a "closed species".</xeblog-conv>

Yes, really. The "Cadey" character you see on the blog was originally envisioned as an [orcadragon](https://www.deviantart.com/junkyardrabbit/journal/Orcadragon-Species-Information-516768513) and I was going to get a reference sheet as such but then I found out about the clusterfuck that is a "closed species". I thought about messing up the details and calling it a drakonorka (lit: dragon-orca in Lojban), but I ended up just cutting off the horns and wings and calling it a day. It made it easier to get commissions that way; a surprising number of artists buy into the "closed species" thing.

<xeblog-sticker name="Cadey" mood="percussive-maintenance"></xeblog-sticker>

I have since decided to pay into the "closed species" ecosystem and I have two "licenses" for orcadragons for both Cadey and another character I've been incubating for the blog named [Palima](https://xeiaso.net/blog/sleeping-the-technical-interview). Once I get the reference sheets, I plan to get the Cadey stickers redrawn as well as bring Palima into the blog properly as the God of Haskell.

<xeblog-conv name="Cadey" mood="percussive-maintenance">How do I prove this? I can't! You just have to trust me, I guess.</xeblog-conv>

### Tying it together with NFTs

<xeblog-hero ai="Epoch4" file="shark-hike" prompt="shark anthro solo female forest hiking walking_stick tshirt denim pants druid long_hair purple_hair red_eyes facing_away"></xeblog-hero>

So with all this in mind, what is an NFT really? If it's not the art associated with that token, what is it? Let's take a look at the word "fungible" a bit closer. When something is fungible, it means that the thing is interchangeable with others like it. A normal white pair of socks is basically identical to another normal white pair of socks. If you have a dollar bill, you can replace it with basically any other dollar bill and the value doesn't really change. A non-fungible item can't be replaced like this. Consider the difference between _any_ pair of white socks and the pair of white socks that you wore every day while you were biking to school during a PhD program. There's value to that pair of socks beyond the inherent value of those socks as a way to cover your body. Non-fungible tokens are not fungible either. There's metadata that lets you say that you have _a specific_ token by ID, not just any token in the collection. It's not a fungible thing in the same way that bread and oil are.

<xeblog-conv name="Cadey" mood="coffee">I am skipping over a _lot_ of economics speak here though. There are people much better suited to talking about this. I am not one of them.</xeblog-conv>

One great thing about humans is that we can _intend_ for things to mean one thing, but they will end up meaning _vastly different things_ to others. NFTs are a great example of this principle in action. Common understanding has owning an NFT to mean that you _own the art_ associated with it. Consider the example of the [people who bought that copy of Dune](https://interestingengineering.com/culture/an-nft-group-bought-a-copy-of-dune-for-304-million-thinking-its-the-copyright) (one of the best science-fiction/fantasy books ever written) and were planning to adapt the book into a movie adaptation with a "decentralized autonomous organization (DAO)" at the reigns.

<xeblog-conv name="Cadey" mood="coffee">I'm not going to get into DAOs in this article. I don't have enough experience in them to speak confidently about what I think about them. I'm not quite sure what purpose they serve though.</xeblog-conv>

There's not been any court cases to set legal precedent for what owning an NFT means. In the eyes of the blockchain people, they seem to think that owning the NFT means you own the copyright to the thing. This could mean that if you make a TV series planned around a character on an image associated with an NFT and then get the NFT social engineered away from you, you [may have to cancel plans for the show unless you pay off the hacker to get the token back](https://www.indiewire.com/2022/06/seth-greene-recovers-lost-nft-1234733180/).

<xeblog-conv name="Open_Skies" mood="flop">Or you'll have to admit that the blockchain isn't the authority on intellectual property and your performative pretending that it is an authority is a scam in order to get money from people.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">"Code is law" my ass.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="wave">Well, to some extent, an NFT can legally mean whatever you write a contract (traditional, not smart) for the NFT to mean. But that would require the court to decide on it, but I think there is an argument to be made that this would fall under the same umbrella as licensing terms, where you can sue someone for violating your licensing terms and require them to pay damages. That’s all speculation though, I haven’t heard of it being tested in court. “An argument to be made” doesn’t really mean much though; once upon a time someone made an argument that all internet moderation is illegal based on a legal precedent that a corporate town couldn’t prevent people from having fliers in a magazine display or something, and therefore we had to let them say anything they wanted anywhere they wanted or we’re violating the law. So we’ll see what happens in court when actual lawyers deal with it.</xeblog-conv>

Without any formal guidance on what all this means or any kind of legal precedent for how to handle all this, it's very unclear what owning an NFT actually does in the eyes of the law. If the NFT goes up in value by a significant amount of money, does that mean you owe capital gains tax? The government has really not caught up with the masses yet and a lot of cryptocurrency technology really exists in a gray area. The kind of gray area that gets dicey with the IRS. You do not want to get dicey with the IRS.

<xeblog-conv name="Cadey" mood="coffee">It seems that the US is treating cryptocurrency accounts like foreign accounts. I'm personally planning on filing my cryptocurrency accounts in my FBAR with my other Canadian accounts come tax season. It's probably overkill, but it's better to give the IRS more information than they need. Filing taxes is going to be _fun_ this year. The sacrifices I make for you people's amusement.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="loaf">I don’t own any cryptocurrency, so I don’t have to deal with the taxes!</xeblog-conv>

I think it is probably best to treat NFTs as something close to sticker or stamp collecting. The tokens themselves are probably inherently worthless (save of course the minimum cost needed to deploy the contract associated with them and the gas fees to mint the tokens), but together they likely have sentimental value to the owner as well as an "aesthetic value" depending on what the creator of that NFT collection was doing when they designed the art itself.

There's also NFTs called [Proof of Attendance Protocol](https://medium.com/new-to-nfts/what-is-a-poap-nft-2a524bb81123) or POAP. These are like ticket stubs at concerts. They let you have collectables associated with specific events, gatherings, or other such things. There's a POAP that's there for anyone who was at the party to celebrate Ethereum's merge to the proof of stake (read: not burning down the forests as much for magic internet money) chain. I personally haven't seen a super good use for these yet, but I can confirm that they are a thing that exists. I really wish I had a good example of how to use these, but I just don't.

<xeblog-conv name="Open_Skies" mood="wave">Something something digital resume fodder to show that you went to conferences/school? IDK.</xeblog-conv>

To be honest though, _if implemented correctly_ NFTs could really help the furry community. Ownership of "adoptables" or the right to use a "closed species" character would be unambiguous. They could even have resale value that gives a commission back to the original creator (provided the relevant smart contracts have a clause to factor this into the equation); and there would be a neutral registry that _anyone_ can read from and confirm ownership of characters in cases where that is relevant to the participants.

<xeblog-conv name="Open_Skies" mood="snug">I could argue that adoptables and closed species and OCs and such are *already* de facto NFTs, just without any automation or engineered resistance to bad actors. Most artists will request the permission from the “owner” of the OCs before they draw the character. Adoptables are sold like NFTs, with an artist creating them whenever they feel like and putting them up for sale in a gallery that users can purchase from and they then “own” that unique art and rights to the character, though the artwork itself is still freely available online for anyone who likes to view or download.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">Or just right-click on the art and put it into the hypercloud. Take that, liberals!</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">So, in practice, furry artwork is *already* NFTs, just implemented as an ad-hoc, undocumented, unverifiable social contract with bad actors, punishing the non-punishers, no way to actually check claims of wrongdoing or see if they’ve been resolved, and multiple incompatible social norms that people can’t just accept federation and say which section they’re in because none of them are documented and that person is the enemy and must be brought into line. Oh, and also beholden to the capricious whims of puritans putting pressure on payment processors or Patreon or whatever to rugpull artists and keep their money. This is *NOT GOOD*.</xeblog-conv>

We can't use any of that in furry-land because the NFT bros pissed off furry artists by scraping twitter, deviantart, and e621 to flip into cheap NFT grifting so much that the entire technology space has become wrongthink to some people. By doing all of this research and writing this article I have probably made some people block me, _even though I'm ending up concluding that all of this space is a bad thing as currently implemented_. It is difficult to describe the level of contempt that a lot of furries have for this technology. I have seen communities with previously irreconcilable differences come together to try and destroy scammers flipping their art into NFTs. So yeah, this entire tech would be _revolutionary_ to the furry community, but nobody wants anything to do with it because idiots pissed in the pool and now it just smells like aging piss.

<xeblog-conv name="Open_Skies" mood="wave">Though I am very glad that furry culture does resist corporatization so strongly, because it would be a horrible loss for the beauty, diversity, and expression of the community to be paved by megacorp lawyers into a metaphorical parking lot or suburban lawn. I’m also very glad that furry culture has such strong social norms for supporting artists and making sure that such critical members of the community that produce so much value actually get to make a living from it. There is less paywalling in furry art than many other kinds of art, while furry artists also get reliably paid more than many other kinds of independent artists, and that’s great! Maybe with some effort we could solve some social and technical problems to reduce paywalls while also increasing the benefit of the artists, and make the environment better for everyone. But that seems… difficult right now. Maybe someday. For now we can just keep upholding the furry values of inclusion, freedom, kindness, weird porn, and punching nazis.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet"><xeblog-picture path="blog/ethereum/furry-fandom-protect"></xeblog-picture><br />Oh yeah, by the way: they want to use NFTs in video games.</xeblog-conv>

<xeblog-conv name="Mara" mood="hmm">What? How does that make any sense?</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Buckle up...</xeblog-conv>

## They want to share items in games

So imagine a world where you can share items between games. Let's say that if you get a weapon in Splatoon you should be able to use it in other games. So you get a [Splatana Wiper](https://splatoonwiki.org/wiki/Splatana_Wiper) from Sheldon and you really like how it plays in the Splatoon meta. Then you close Splatoon and open up Call of Duty. The NFT bro dream is that you'd be able to pull the Splatana out of Splatoon and then drop it into your Call of Duty loadout so that you can whack people in the face with it for 120 damage.

Imagine if this could spread to other items too. Unlocking an outfit or weapon class for Sena in Xenoblade Chronicles 3 means you can take that outfit to other games like Fortnite. Grab a bunch of items from Harmony in Splatsville and then use them all to decorate your house in Animal Crossing!

<xeblog-conv name="Mara" mood="hmm">How would any of this actually...work? This sounds nice yeah but how would it make any sense? Maybe sharing items between Animal Crossing and Splatoon makes sense (they share a development team and game engine even). But Call of Duty and Splatoon are about as polar opposites as Quake and Chess. Wouldn't that look like a mess kinda like VRChat or Second Life?</xeblog-conv>

This won't work at all in practice. Every one of those games I listed has its own unique art style and that artstyle dictates a lot about the vision and direction of the games in question. 

<xeblog-conv name="Cadey" mood="coffee">Okay, Fortnite is a bit of an exception. It _used to_ have a cohesive artstyle at the beginning, but over time it's turned into a game where Rick from Rick and Morty, Goku, Naruto, and [Marshmello](https://en.wikipedia.org/wiki/Marshmello) can get into a squad together and Hadouken people to death; it's more surreal than it sounds.</xeblog-conv>

This would also make game balance nearly impossible. Taking that example of moving a Splatana from Splatoon into Call of Duty more literally, the Splatana does 30 damage when the fire button is tapped and up to 120 damage when the fire button is held and you beat someone in the face with it. This is balanced because each squid/octopus child thing has 100 hitpoints, so 4 basic shots splats an opponent and if you get good at reading people you can oneshot people around corners. The base player health in Call of Duty is anywhere from [100 to 150](https://callofduty.fandom.com/wiki/Health_System). Depending on the game this could mean that the Splatana can either be about normally powered or hilariously underpowered to the point that you wouldn't want to use it anyways.


<xeblog-conv name="Open_Skies" mood="snug">And then you bring it into Minecraft where players have 20 hit points and every attack is a one hit kill. But wait, do you lose it across all games when you die? Can you bring it with you every time you spawn?</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">Does the ink coverage from the Splatana stay in the equation? Do you also get the Ultra Stamp special? How would games know to account for this?</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="loaf">In Planetary Annihilation and Supreme Commander, the ink coverage from the Splatoon weapons would obviously be an auxiliary power generation mechanic, representing some kind of nanomachine based solar array. Which means that the meta for competitive play would require purchasing and progressing to high levels in games from completely different genres, and I can’t imagine the gaming community putting up with that. It would also require quadratic-or-worse balancing efforts where every game needs to care about the balance of every other game and write code to handle bringing items and mechanics in from every other game, and/or it would allow anyone who feels like self-publishing a game to home-brew a balance breaking super item in whatever game they want, and I can’t imagine that game developers would put up with that.<br /><br />Oh, you missed something. Tie in mobile games released alongside a main desktop game occasionally share items between the two.</xeblog-conv>

We do have one thing that is close to this ideal, but it's very limited as you'd imagine. VTuber software is segmented into two types, 2D and 3D. I can't speak much about the 2D software, but for 3D VTubers we use a model format called [VRM](https://vrm.dev/en/) as a vendor-neutral interchange format. I can pull my VRM avatar into other games such as [Synth Riders](https://synthridersvr.com/) and dance around on stream:

<center><iframe width="560" height="315" src="https://www.youtube-nocookie.com/embed/tC0tgvHqHqw" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe></center>

I can take that same model and plunk it into Godot and then animate it out to do whatever I want. The artstyle is maintained because all of these avatars have shading metadata that denotes them as anime styled. This is mainly because the VRM format was made by Japanese VTuber software developers and mostly because that's what the whole project was intended to be used for: anime styled VTuber characters.

<xeblog-conv name="Cadey" mood="facepalm">It's pretty nice in practice, but there's a depressing amount of software out there such as VRChat that won't let me import a pre-rigged model, which means that I have to do a [variety of annoying hacks](https://xeiaso.net/blog/vrchat-avatar-to-vrm-vtubing-2022-01-02) to get things working.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="wave">I have a few disagreements with the design of VRM, and a lot of annoyances at the lack of support for it. But that discussion is for another time.</xeblog-conv>

<xeblog-conv name="Mara" mood="happy">At least the VRM format works and it's already supported by multiple engines. It will get adoption in time, it can also be extended if that is really needed. I really wish there were more games supporting it.</xeblog-conv>

## What is Ethereum worth?

<xeblog-hero ai="Waifu Diffusion v1.2" file="botw-space-needle" prompt=" the legend of zelda breath of the wild, mountain, mid-day, studio ghibli, space needle, river, guardian turret"></xeblog-hero>

I've been talking about Ethereum for most of this article because Ethereum is the cryptocurrency that everyone seems to base all their work on. It was the first chain where smart contracts and NFTs really took off. However, I've never really covered what Ethereum is _really worth_. What is an Ethereum token? How does it have value?

<xeblog-conv name="Numa" mood="delet"><br />&gt;money<br />&gt;inherent value<br /><br />pick one.</xeblog-conv>

<xeblog-conv name="Cadey" mood="facepalm">Okay, yes, from the philosophical standpoint modern currencies don't really have value backed by anything. Once the US abolished the [Gold standard](https://en.wikipedia.org/wiki/Gold_standard) the value of the US dollar was detached from the value of physical hunks of gold and it became a [fiat currency](https://en.wikipedia.org/wiki/Fiat_money). So in essence a dollar bill in your pocket doesn't really have its own value outside of the mass psychosis that everyone experiences to think that the US dollar has value. At the same time though, the US government has this thing called "an army", which kind of enforces the idea that the US dollar has value.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">The inherent value of one US dollar is that you could use it as toilet paper or [burn it for heat in the winter](https://library.randolphschool.net/c.php?g=237930&p=1581974) or [wrap something in it](https://www.lightspeedmagazine.com/fiction/the-cambist-and-lord-iron-a-fairy-tale-of-economics). Everything else comes from that little bit of text on it “This note is legal tender for all debts public and private” and the government enforcing it and the society humoring the collective delusion.</xeblog-conv>

<xeblog-picture path="blog/ethereum/money_fireplace"></xeblog-picture>

<xeblog-conv name="Numa" mood="delet">Something something monopoly on violence something something.</xeblog-conv>

So with all this in mind, what is Ethereum really worth? Honestly after doing a bit of research I can't really find a solid answer to this seemingly simple question. Money is complicated. Economics is a very nihilist field of study when you peel back all the layers of terminology and math.

Ethereum used to be a proof of work consensus blockchain, which is a bunch of buzzwords that means that some nodes on the network called "miners" would have GPUs solving sudokus 24/7 in order to get a hash that met certain criteria. When that happened, the miner would be awarded with a block's worth of tokens. Ethereum is also unlike other cryptocurrencies because Ethereum has no supply cap (the maximum number of tokens that can exist) unlike Bitcoin that has a cap of 21 million bitcoins ever being able to exist.

So if Ethereum is no longer a proof of work chain (IE: the price of an Ethereum token had value from the power wasted by the miner attempting to create it) and there's no cap of how many tokens can exist, what is Ethereum really worth?

<xeblog-conv name="Open_Skies" mood="snug">The difficulty of the creation doesn’t create value, it creates scarcity. The value needs to come from something else. This is where gas fees come in. It’s tempting to say that the value of Ethereum is the amount of computation that can be purchased on the Ethereum blockchain by burning the eth for gas. But that doesn’t explain the price; computing on the Ethereum blockchain is really inefficient because of all the validators double checking the computations, so the amount of money it would cost to buy the computation on a cloud compute provider and the amount it would cost to buy on Ethereum if nothing else were in play are very different, and the actual price of gas doesn’t match that either. So they have to be paying for *something*, and that something is legitimacy. We can view Ethereum as a machine for taking one big agreement, the agreement about the consensus algorithm, the smart contract instruction set and runtime, and what the current head of the chain is, and transforming that into a billion tiny agreements about who owns what and how to handle every little transaction. In this view, the value of gas, and by extension the value of Ethereum, is the value of making trustworthy agreements about digital objects. And it turns out that globally accessible and verifiable trustworthy agreements about digital objects are *extremely valuable*.</xeblog-conv>

<xeblog-conv name="Mara" mood="hmm">This is something that we haven't really considered before. Before that argument I thought that Ethereum really just had the same worth that you get with government issued currency: everyone just agrees it has value and trading volume along with "market confidence" would dictate the value. Ethereum being worth the legitimacy of the system is a vastly different thing to think about.<br /><br />What about the Ethereum VM? Don't you use Ethereum as "gas" for the virtual machine that processes all the smart contract executions and transactions? How does that fit into it all?</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="wave">Not quite. You buy gas with Ethereum. The Ethereum VM is the language by which you describe the rules of how something is legitimate. If you want to make something that has legitimacy, you need to explain exactly how legitimate interactions with it go. For example, a legitimate interaction with a store is you pay a specific price and in exchange receive a specific good; if you do not offer enough, your money isn’t taken and you don’t get the good; if you offer too much you get the good and any excess back. A legitimate interaction with a voting booth is that you mark down your choice on each part of the ballot, possibly including marking abstain fields, and submit it to the vote exactly once; the voting booth guarantees that the ballot you cast will be included in the results and that no fictitious ballots that no person submitted or multiple ballots from the same person are included. These rules have to be written down to be enforced, and the more complicated the rules are to check and more states they depend on the more difficult they are to enforce: thus, the Ethereum VM is how the rules of legitimacy are written down so everyone can agree on them, and gas usage is the measure of how difficult it is to perform the transaction according to the rules so everyone agrees on it. When you submit a transaction to the Ethereum network, you include in it a marker for how much you are willing to pay per unit of gas, and how much gas you are willing to buy before you would rather not have the transaction go through. Different miners might have different amounts that they’re willing to sell gas for, and so being willing to pay more for gas means that a miner that is willing to accept that rate will be selected sooner and the transaction will complete sooner, and if you don’t pay enough a miner will never accept it. So this makes it so that you can effectively pay a premium for expedited processing of your transactions or pay lower rates for slower inclusions. Ethereum is the only payment you can buy gas with, so apart from regarding Ethereum as a speculative asset, gas is really where the value gets meaning. It’s the atomic cost of ensuring legitimacy for a single step of operation. (Things can be legitimate without the legitimacy being checked; it’s possible that the weird person on the subway with a deed written on toilet paper in crayon offering to sell the Mona Lisa for $100 actually is the legitimate owner and will honor the deal, but in order for other people to actually trust the legitimacy it needs to be checked in some way that people can trust, which takes effort and coordination and thus has to cost *something*, even if it’s "a miniscule amount of credits from a UBI or publicly funded service" or "a nebulous tiny expenditure of social capital among the community" (if you’re wondering what something nebulous like that costs, just imagine that someone used and abused the thing so much that the providers had to stop them; what did they run out of that caused them to get dropped))</xeblog-conv>

<xeblog-conv name="Mara" mood="hmm">That is a lot to think about. I don't understand money at all.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">Oh just wait until we get into the unrestrained fun that is _user-defined_ fungible tokens on the blockchain.</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">That may have to wait for part 3. This article is already an hour long.</xeblog-conv>

I think the main takeaway here is that Ethereum is _fundamentally worthless_ in the same way that US dollars are, but fiat mass psychosis and trading volume gives it value.

<xeblog-conv name="Open_Skies" mood="wave">Could you explain a bit more about what you mean about the value coming from mass psychosis? I’m not sure I agree.</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Mass psychosis is probably a bit harsh, but I'm trying to get at that shared cultural delusion that money has value because everyone thinks it does. Like, what is a dollar really worth? What does it really signify?</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="wave">~26 grams of silver~</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">Not anymore!</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">Yeah, now a dollar is worth 1.4 grams of silver. A dollar is worth what you can buy for it. You can sell a dollar for anything in the dollar store. It's also worth a bit under a pound of rice.</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">This is going to end with financial nihilism even more isn't it.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="idle">I'm kinda already there? And from my perspective of financial nihilism, calling it mass psychosis is misleading; the value of a dollar is not a mass delusion, it is a reasonable, evidence based belief that you will be able to exchange the dollar for a given amount of stuff later. It's just that there is no centrally mandated exchange rate to some other currency. Because gold is just another currency; you can't eat gold, and only recently did it become useful in industrial processes that the average person cares about/is materially affected by (electronics manufacture, among others (which isn't something the average person can do on their own, so that doesn't affect the argument)), so gold is acting as just another currency: it's value is based on the reasonable, evidence based belief that you will be able to trade it for a given amount of the stuff you want in the future. And it isn't pegged to anything, so the value can fluctuate significantly as the market changes. So, way back when, a bank note/IOU/check from any local bank, an appraiser's note, a dollar, a physical lump of gold, etc. were all fiscally equivalent in some sense by being denominated in terms of the same thing. It would be like making a currency in the modern day where each note of it corresponds to a given number of USD. (Like say $100 USD notes, which are pegged to the price of $1 USD notes, since you can exchange one for the other at any bank (though the pegging isn't perfect; there is some fluctuation in the value so that the values of 1 $100 USD note and of 100 $1 USD notes aren't exactly the same) or just a stablecoin.)</xeblog-conv>

<xeblog-conv name="Mara" mood="hmm">A stable-what?</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">It's definitely something for at least part 3. "Stablecoins" are a topic and a half.</xeblog-conv>

## Capitalism and ownership

<xeblog-hero ai="Waifu Diffusion v1.2" file="bloodborne-bad-end" prompt="populated bloodborne old valley with a obscure person at the centre and a ruined gothic city in the background, trees and stars in the background, falling red petals, epic red - orange moonlight, perfect lightning, wallpaper illustration by niko delort and kentaro miura, 4k"></xeblog-hero>

This all leaves me with a bit of a sour note. The idea that they are peddling is that cryptocurrency is a means of liberation against a "broken" financial system, but then the solution is just another means of limitation. Want to send a payment to a friend? That will cost ten cents to a dollar depending on the time of day. Want to get a new NFT? Ten cents to a dollar. Want to be in this club? Better get one of the tokens! They're all sold out? Too bad! Buy one second-hand! It's really just another method of dividing people into haves and have-nots. God knows we need more of them.

All of this really just adds more and more layers to this game of capitalism that we are all forced to play against our wills. Sure, the model of cryptocurrency allows you _relative freedom_ within its rules, but there are still fundamental limitations to what you can do with things outside its cryptosystem. One of the artists I like to commission lives in Russia, so paying them is difficult (especially after Paypal pulled out of there) except for when I use cryptocurrency. 

Don't get me wrong, cryptocurrencies work great when the system is self-contained, but as soon as real money and assets get involved it becomes a huge mess. This relative freedom is another layer of capitalistic hierarchy and it helps the early adopters and the already rich get richer (through transaction fees and block rewards). 

In practice though, the distributed ledger usually ends up having a few big players that sell blockchain API access as a service. The blockchain for Ethereum is over one terabyte. This is an _insane_ amount of disk space. There are efforts to lessen this, but at the time of writing the blockchain is just _too big_ for most people to use themselves. So you need blockchain access as a service. Which means you need to pay up to use your distributed money in the first place. This kind of defeats the point of having decentralized money if people are just going to recentralize around large players to be able to use it.

<xeblog-conv name="Open_Skies" mood="flop">One terrifying consequence of the way people interact with Ethereum is the existence of the generalized transaction frontrunning bot. We’ll have to get into what that actually means in a future part, but here there be monsters.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">The masses think they are free to send money to each other but the rich just get richer off their backs along the way. As above, so below.</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">This whole thing is a mess. I only play this capitalism game because I have no other choice, not because I want to.</xeblog-conv>

<xeblog-conv name="Open_Skies" mood="flop">Yeah. I just want to make useful things and help people, and only being allowed to make useful things at scale if I also make them artificially less useful in order to ensure I capture the value created by them is annoying.</xeblog-conv>

## Conclusion

I don't really have a good conclusion here. I don't know how I feel about Ethereum. It's got me split down the middle because there are _good ideas_ here. There are things that if done correctly are genuinely useful. There's just so many grifters. I didn't want to go into the grifters too much in this article, there is going to be more in future parts to this series. 

While I was writing this article I was contacted by someone representing a "play to earn" game (one where you get NFTs that act as game objects for playing a game connected to the blockchain) offering me 0.55ETH (about USD$731 at the time of writing) to make three tweets shilling their game that had _22 online players_ and no public client download. $731 isn't even worth the tax burden (including my having to explain things to my accountant, who is probably going to go on a furious googling spree), not to mention the reputational damage that I would incur for shilling such a "product".

<xeblog-conv name="Cadey" mood="coffee">Perfect timing, eh?</xeblog-conv>

All of this technology is really cool, I just wish I could use it for anything productive. I'm able to commission that one artist with it at least. Maybe there will end up being some kind of killer app with all of this, but until then I'm just gonna keep playing video games and writing these near novel-length articles. I'm also going to have to continue blocking and reporting the scams people shovel at me on Twitter.

<xeblog-hero ai="Waifu Diffusion v1.2" file="blockchain-hacker-miku" prompt="hacker hoodie 1girl Ethereum blockchain hatsune_miku twintails blue_hair laptop coffee starbucks blue_eyes anime waifu super hacker very cute"></xeblog-conv>

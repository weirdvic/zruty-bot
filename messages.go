package main

const (
	greetAdmin     string = `Приветствую!`
	notAdmin       string = `Вы не являетесь администратором этого бота.`
	kickMessage    string = `Пользователь %s покидает чат: %s`
	welcomeMessage string = `Добро пожаловать, <a href="tg://user?id=%d">%s</a>!
Этот чат посвящён компьютерной игре NetHack.
Полезные ссылки есть в информации о группе.
Войдя в чат, вы ОБЯЗАНЫ поздороваться с участниками.
На это у вас есть 23 часа.
В противном случае вы будете забанены.`
)

var deathCauses []string = []string{
	"ascended",
	"burned by a tower of flame",
	"burned by molten lava",
	"caught herself in her own ball of cold",
	"caught herself in her own fireball",
	"caught himself in his own ball of lightning",
	"caught himself in his own burning oil",
	"caught himself in his own death field",
	"caught himself in his own magical blast",
	"choked",
	"committed suicide",
	"crunched in the head by a iron ball",
	"crushed to death by a collapsing drawbridge",
	"crushed to death by a exploding drawbridge",
	"crushed to death by a falling drawbridge",
	"died of starvation",
	"dissolved in molten lava",
	"dragged downstairs by a iron ball",
	"drowned in a moat of water",
	"drowned in a pool of water by a kraken",
	"drowned in a pool of water by a python",
	"drowned in deep water",
	"escaped (in celestial disgrace)",
	"escaped with a fake Amulet",
	"fell into a pit of iron spikes",
	"fell into a pit",
	"fell onto a sink",
	"killed by Ashikaga Takauji",
	"killed by Asmodeus",
	"killed by Croesus",
	"killed by Death",
	"killed by Demogorgon",
	"killed by Dispater",
	"killed by Famine",
	"killed by Ixoth",
	"killed by Izchak the shopkeeper",
	"killed by Juiblex",
	"killed by Lord Surtur",
	"killed by Master Kaen",
	"killed by Nalzok",
	"killed by Orcus",
	"killed by Pestilence",
	"killed by Vlad the Impaler",
	"killed by Yeenoghu",
	"killed by a Green-elf",
	"killed by a Grey-elf",
	"killed by a Keystone Kop",
	"killed by a Kop Kaptain",
	"killed by a Kop Lieutenant",
	"killed by a Kop Sergeant",
	"killed by a Mordor orc",
	"killed by a Woodland-elf",
	"killed by a acid venom",
	"killed by a air elemental of Tyr",
	"killed by a ape",
	"killed by a arrow",
	"killed by a baby black dragon",
	"killed by a baby blue dragon",
	"killed by a baby crocodile",
	"killed by a baby green dragon",
	"killed by a baby orange dragon",
	"killed by a baby red dragon",
	"killed by a baby silver dragon",
	"killed by a baby white dragon",
	"killed by a balrog",
	"killed by a baluchitherium",
	"killed by a barbed devil",
	"killed by a barrow wight",
	"killed by a bat",
	"killed by a bear trap",
	"killed by a black dragon",
	"killed by a black naga hatchling",
	"killed by a black naga",
	"killed by a black pudding",
	"killed by a black unicorn",
	"killed by a blast of acid",
	"killed by a blast of disintegration",
	"killed by a blast of fire",
	"killed by a blast of frost",
	"killed by a blast of lightning",
	"killed by a blast of missiles",
	"killed by a blue dragon",
	"killed by a blue jelly",
	"killed by a boiling potion",
	"killed by a bolt of cold",
	"killed by a bolt of fire",
	"killed by a bolt of lightning",
	"killed by a bone devil",
	"killed by a boulder",
	"killed by a brown mold",
	"killed by a bugbear",
	"killed by a burning book",
	"killed by a burning potion of oil",
	"killed by a burning scroll",
	"killed by a cadaver",
	"killed by a captain",
	"killed by a carnivorous ape",
	"killed by a carnivorous bag",
	"killed by a cave spider",
	"killed by a centipede",
	"killed by a chameleon imitating a air elemental",
	"killed by a chameleon imitating a arch-lich",
	"killed by a chameleon imitating a baluchitherium",
	"killed by a chameleon imitating a hell hound",
	"killed by a chameleon imitating a jabberwock",
	"killed by a chameleon imitating a leocrotta",
	"killed by a chameleon imitating a master lich",
	"killed by a chameleon imitating a mastodon",
	"killed by a chameleon imitating a minotaur",
	"killed by a chameleon imitating a vampire lord",
	"killed by a chameleon imitating a vampire",
	"killed by a chameleon imitating a winged gargoyle",
	"killed by a chickatrice",
	"killed by a clay golem",
	"killed by a cloud of poison gas",
	"killed by a cobra",
	"killed by a cockatrice",
	"killed by a contact-poisoned spellbook",
	"killed by a contaminated potion",
	"killed by a couatl of The Lady",
	"killed by a couatl",
	"killed by a coyote",
	"killed by a crocodile",
	"killed by a crossbow bolt",
	"killed by a dagger",
	"killed by a dart",
	"killed by a death ray",
	"killed by a demilich",
	"killed by a dingo",
	"killed by a disenchanter",
	"killed by a djinni",
	"killed by a dog",
	"killed by a dust vortex",
	"killed by a dwarf king",
	"killed by a dwarf lord",
	"killed by a dwarf mummy",
	"killed by a dwarf zombie",
	"killed by a dwarf",
	"killed by a dwarvish spear",
	"killed by a earth elemental of Lugh",
	"killed by a exploding large box",
	"killed by a exploding potion",
	"killed by a fall onto poison spikes",
	"killed by a falling object",
	"killed by a falling rock",
	"killed by a fire ant",
	"killed by a fire elemental of Kos",
	"killed by a fire elemental",
	"killed by a fire giant",
	"killed by a fire vortex",
	"killed by a flaming sphere",
	"killed by a flesh golem",
	"killed by a fog cloud",
	"killed by a forest centaur",
	"killed by a fox",
	"killed by a freezing sphere",
	"killed by a frost giant",
	"killed by a gargoyle",
	"killed by a garter snake",
	"killed by a gas cloud",
	"killed by a gas spore's explosion",
	"killed by a gecko",
	"killed by a gelatinous cube",
	"killed by a ghost",
	"killed by a ghoul",
	"killed by a giant ant",
	"killed by a giant bat",
	"killed by a giant beetle",
	"killed by a giant mimic",
	"killed by a giant mummy",
	"killed by a giant rat",
	"killed by a giant spider",
	"killed by a giant zombie",
	"killed by a glass golem",
	"killed by a glass piercer",
	"killed by a gnome king",
	"killed by a gnome lord",
	"killed by a gnome mummy",
	"killed by a gnome zombie",
	"killed by a gnome",
	"killed by a gnomish wizard",
	"killed by a goblin",
	"killed by a gold golem",
	"killed by a golden naga hatchling",
	"killed by a golden naga",
	"killed by a gray dragon",
	"killed by a gray ooze",
	"killed by a gray unicorn",
	"killed by a green dragon",
	"killed by a green mold",
	"killed by a gremlin",
	"killed by a grid bug",
	"killed by a guard",
	"killed by a guardian naga hatchling",
	"killed by a hell hound pup",
	"killed by a hell hound",
	"killed by a hezrou",
	"killed by a hill orc",
	"killed by a hobbit",
	"killed by a hobgoblin",
	"killed by a homunculus",
	"killed by a horned devil",
	"killed by a horse",
	"killed by a housecat",
	"killed by a human mummy",
	"killed by a human zombie",
	"killed by a iguana",
	"killed by a incubus of Ishtar",
	"killed by a iron ball collision",
	"killed by a iron golem",
	"killed by a jabberwock",
	"killed by a jackal",
	"killed by a jaguar",
	"killed by a ki-rin",
	"killed by a killer bee",
	"killed by a kitten",
	"killed by a kobold lord",
	"killed by a kobold mummy",
	"killed by a kobold shaman",
	"killed by a kobold zombie",
	"killed by a kobold",
	"killed by a land mine",
	"killed by a large cat",
	"killed by a large dog",
	"killed by a large kobold",
	"killed by a large mimic",
	"killed by a leather golem",
	"killed by a leocrotta",
	"killed by a leprechaun",
	"killed by a lich",
	"killed by a lieutenant",
	"killed by a little dart",
	"killed by a little dog",
	"killed by a lizard",
	"killed by a long worm",
	"killed by a lurker above",
	"killed by a lynx",
	"killed by a magic missile",
	"killed by a magical explosion",
	"killed by a manes",
	"killed by a marilith",
	"killed by a master lich",
	"killed by a master mind flayer",
	"killed by a mastodon",
	"killed by a mildly contaminated potion",
	"killed by a mind flayer",
	"killed by a minotaur",
	"killed by a mountain centaur",
	"killed by a mumak",
	"killed by a nalfeshnee",
	"killed by a newt",
	"killed by a ninja",
	"killed by a nurse",
	"killed by a ogre king",
	"killed by a orange dragon",
	"killed by a orc shaman",
	"killed by a owlbear",
	"killed by a panther",
	"killed by a paper golem",
	"killed by a partisan",
	"killed by a pit fiend",
	"killed by a pit viper",
	"killed by a plains centaur",
	"killed by a poison dart",
	"killed by a poisoned blast",
	"killed by a poisoned needle",
	"killed by a poisonous corpse",
	"killed by a pony",
	"killed by a potion of acid",
	"killed by a potion of holy water",
	"killed by a potion of unholy water",
	"killed by a priestess",
	"killed by a psychic blast",
	"killed by a purple worm",
	"killed by a pyrolisk",
	"killed by a python",
	"killed by a quivering blob",
	"killed by a rabid rat",
	"killed by a raven",
	"killed by a red dragon",
	"killed by a red mold",
	"killed by a red naga hatchling",
	"killed by a red naga",
	"killed by a riding accident",
	"killed by a rock mole",
	"killed by a rock piercer",
	"killed by a rock troll",
	"killed by a rock",
	"killed by a rope golem",
	"killed by a rothe",
	"killed by a salamander",
	"killed by a sasquatch",
	"killed by a scorpion",
	"killed by a scroll of earth",
	"killed by a scroll of genocide",
	"killed by a sergeant",
	"killed by a sewer rat",
	"killed by a shade",
	"killed by a shattered potion",
	"killed by a shocking sphere",
	"killed by a shopkeeper",
	"killed by a shuriken",
	"killed by a silver dragon",
	"killed by a skeleton",
	"killed by a small mimic",
	"killed by a snake",
	"killed by a soldier ant",
	"killed by a soldier",
	"killed by a spear",
	"killed by a spotted jelly",
	"killed by a stalker",
	"killed by a steam vortex",
	"killed by a stone giant",
	"killed by a stone golem",
	"killed by a storm giant",
	"killed by a straw golem",
	"killed by a succubus of Venus",
	"killed by a succubus",
	"killed by a system shock",
	"killed by a tengu",
	"killed by a thrown potion",
	"killed by a tiger",
	"killed by a titan",
	"killed by a titanothere",
	"killed by a touch of death",
	"killed by a tower of flame",
	"killed by a trapper",
	"killed by a troll",
	"killed by a umber hulk",
	"killed by a unsuccessful polymorph",
	"killed by a vampire bat",
	"killed by a vampire in bat form",
	"killed by a vampire lord in bat form",
	"killed by a vampire lord",
	"killed by a vampire",
	"killed by a violet fungus",
	"killed by a vrock",
	"killed by a wand",
	"killed by a warg",
	"killed by a warhorse",
	"killed by a watch captain",
	"killed by a watchman",
	"killed by a water demon",
	"killed by a water elemental of Mog",
	"killed by a water elemental",
	"killed by a water moccasin",
	"killed by a water troll",
	"killed by a werejackal",
	"killed by a wererat",
	"killed by a werewolf",
	"killed by a white dragon",
	"killed by a white unicorn",
	"killed by a winged gargoyle",
	"killed by a winter wolf cub",
	"killed by a winter wolf",
	"killed by a wolf",
	"killed by a wood golem",
	"killed by a worthless piece of glass",
	"killed by a wraith",
	"killed by a wumpus",
	"killed by a xorn",
	"killed by a zruty",
	"killed by an Aleax of Odin",
	"killed by an Aleax",
	"killed by an Angel of Crom",
	"killed by an Angel",
	"killed by an Archon",
	"killed by an Elvenking",
	"killed by an Olog-hai",
	"killed by an Uruk-hai",
	"killed by an acid blob",
	"killed by an acidic corpse",
	"killed by an acidic glob",
	"killed by an air elemental",
	"killed by an anti-magic implosion",
	"killed by an arch-lich",
	"killed by an earth elemental",
	"killed by an eel",
	"killed by an electric chair",
	"killed by an electric shock",
	"killed by an elf mummy",
	"killed by an elf zombie",
	"killed by an elf-lord",
	"killed by an elven arrow",
	"killed by an elven dagger",
	"killed by an energy vortex",
	"killed by an ettin mummy",
	"killed by an ettin zombie",
	"killed by an exploding chest",
	"killed by an exploding crystal ball",
	"killed by an exploding wand",
	"killed by an explosion",
	"killed by an ice troll",
	"killed by an ice vortex",
	"killed by an imp",
	"killed by an incubus",
	"killed by an iron piercer",
	"killed by an ochre jelly",
	"killed by an ogre lord",
	"killed by an ogre",
	"killed by an orc mummy",
	"killed by an orc zombie",
	"killed by an orc-captain",
	"killed by an orcish arrow",
	"killed by an orcish dagger",
	"killed by an quantum mechanic",
	"killed by an quasit",
	"killed by an queen bee",
	"killed by an yellow dragon",
	"killed by an yeti",
	"killed by axing a hard object",
	"killed by boiling potions",
	"killed by brainlessness",
	"killed by bumping into a wall",
	"killed by burning scrolls",
	"killed by colliding with the ceiling",
	"killed by contaminated water",
	"killed by exhaustion",
	"killed by falling downstairs",
	"killed by genocidal confusion",
	"killed by kicking sink",
	"killed by life drainage",
	"killed by overexertion",
	"killed by petrification",
	"killed by shattered potions",
	"killed by sipping boiling water",
	"killed by strangulation",
	"killed by the Chromatic Dragon",
	"killed by the Cyclops",
	"killed by the Dark One",
	"killed by the Minion of Huhetotl",
	"killed by the Oracle",
	"killed by the Wizard of Yendor",
	"killed by the high priestess of Susanowo",
	"killed by the wrath of a god",
	"killed by touching Mjollnir",
	"killed by touching The Eye of the Aethiopica",
	"killed by touching The Sceptre of Might",
	"killed by tumbling down a flight of stairs",
	"killed by using a magical horn on himself",
	"killed herself by breaking a wand",
	"killed herself with her bullwhip",
	"killed while stuck in creature form",
	"petrified by Medusa",
	"petrified by a captain",
	"petrified by a chickatrice corpse",
	"petrified by a chickatrice egg",
	"petrified by a chickatrice",
	"petrified by a cockatrice corpse",
	"petrified by a cockatrice egg",
	"petrified by a cockatrice",
	"petrified by a lieutenant",
	"petrified by a sergeant",
	"petrified by a soldier",
	"petrified by kicking a chickatrice corpse without boots",
	"petrified by kicking a cockatrice corpse without boots",
	"petrified by losing gloves while wielding a chickatrice corpse",
	"petrified by losing gloves while wielding a cockatrice corpse",
	"petrified by tasting Medusa meat",
	"petrified by tasting chickatrice meat",
	"petrified by tasting cockatrice meat",
	"petrified by touching a chickatrice corpse bare-handed",
	"petrified by touching a cockatrice corpse bare-handed",
	"petrified by trying to tin a chickatrice without gloves",
	"petrified by trying to tin a cockatrice without gloves",
	"poisoned by Demogorgon",
	"poisoned by Juiblex",
	"poisoned by Pestilence",
	"poisoned by Scorpius",
	"poisoned by a centipede",
	"poisoned by a crossbow bolt",
	"poisoned by a dart",
	"poisoned by a fall onto poison spikes",
	"poisoned by a giant spider",
	"poisoned by a killer bee",
	"poisoned by a little dart",
	"poisoned by a pit viper",
	"poisoned by a poison dart",
	"poisoned by a quasit",
	"poisoned by a rabid rat",
	"poisoned by a rotted corpse",
	"poisoned by a rotted glob of gray ooze",
	"poisoned by a scorpion",
	"poisoned by a snake",
	"poisoned by a soldier ant",
	"poisoned by a vampire bat",
	"poisoned by a water moccasin",
	"poisoned by an arrow",
	"poisoned by an orcish arrow",
	"poisoned by an unicorn horn",
	"shot herself with a death ray",
	"slipped while mounting a saddled pet",
	"squished under a boulder",
	"starvation",
	"trickery",
	"turned into green slime",
	"turned to slime by a green slime",
	"went to heaven prematurely",
	"zapped herself with a spell",
	"zapped himself with a wand",
}

package model

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/robertjanetzko/LegendsBrowser2/backend/util"
)

func (x *Honor) Requirement() string {
	var list []string
	if x.RequiresAnyMeleeOrRangedSkill {
		list = append(list, "attaining sufficent skill with a weapon or technique")
	}
	if x.RequiredSkill != HonorRequiredSkill_Unknown {
		list = append(list, "attaining enough skill with the "+x.RequiredSkill.String())
	}
	if x.RequiredBattles == 1 {
		list = append(list, "serving in combat")
	}
	if x.RequiredBattles > 1 {
		list = append(list, fmt.Sprintf("participating in %d battles", x.RequiredBattles))
	}
	if x.RequiredYears >= 1 {
		list = append(list, fmt.Sprintf("%d years of membership", x.RequiredYears))
	}
	if x.RequiredKills >= 1 {
		list = append(list, fmt.Sprintf("slaying %d enemies", x.RequiredKills))
	}

	return " after " + andList(list)
}

func (x *HistoricalEventAddHfEntityHonor) Html() string {
	e := world.Entities[x.EntityId]
	h := e.Honor[x.HonorId]
	return fmt.Sprintf("%s received the title %s of %s%s", hf(x.Hfid), h.Name(), entity(x.EntityId), h.Requirement())
}

func (x *HistoricalEventAddHfEntityLink) Html() string {
	h := hf(x.Hfid)
	c := entity(x.CivId)
	if x.AppointerHfid != -1 {
		c += fmt.Sprintf(", appointed by %s", hf(x.AppointerHfid))
	}
	switch x.Link {
	case HistoricalEventAddHfEntityLinkLink_Enemy:
		return h + " became an enemy of " + c
	case HistoricalEventAddHfEntityLinkLink_Member:
		return h + " became a member of " + c
	case HistoricalEventAddHfEntityLinkLink_Position:
		return h + " became " + world.Entities[x.CivId].Position(x.PositionId).GenderName(world.HistoricalFigures[x.Hfid]) + " of " + c
	case HistoricalEventAddHfEntityLinkLink_Prisoner:
		return h + " was imprisoned by " + c
	case HistoricalEventAddHfEntityLinkLink_Slave:
		return h + " was enslaved by " + c
	case HistoricalEventAddHfEntityLinkLink_Squad:
		return h + " became a hearthperson/solder of  " + c // TODO
	}
	return h + " became SOMETHING of " + c
}

func (x *HistoricalEventAddHfHfLink) Html() string {
	h := hf(x.Hfid)
	t := hf(x.HfidTarget)
	switch x.LinkType {
	case HistoricalEventAddHfHfLinkLinkType_Apprentice:
		return h + " became the master of " + t
	case HistoricalEventAddHfHfLinkLinkType_Deity:
		return h + " began worshipping " + t
	case HistoricalEventAddHfHfLinkLinkType_FormerMaster:
		return h + " ceased being the apprentice of " + t
	case HistoricalEventAddHfHfLinkLinkType_Lover:
		return h + " became romantically involved with " + t
	case HistoricalEventAddHfHfLinkLinkType_Master:
		return h + " began an apprenticeship under " + t
	case HistoricalEventAddHfHfLinkLinkType_PetOwner:
		return h + " became the owner of " + t
	case HistoricalEventAddHfHfLinkLinkType_Prisoner:
		return h + " imprisoned " + t
	case HistoricalEventAddHfHfLinkLinkType_Spouse:
		return h + " married " + t
	default:
		return h + " LINKED TO " + t
	}
}

func (x *HistoricalEventAddHfSiteLink) Html() string {
	h := hf(x.Histfig)
	c := ""
	if x.Civ != -1 {
		c = " of " + entity(x.Civ)
	}
	b := ""
	if x.Structure != -1 {
		b = " " + structure(x.SiteId, x.Structure)
	}
	s := site(x.SiteId, "in")
	switch x.LinkType {
	case HistoricalEventAddHfSiteLinkLinkType_HomeSiteAbstractBuilding:
		return h + " took up residence in " + b + c + " " + s
	case HistoricalEventAddHfSiteLinkLinkType_Occupation:
		return h + " started working at " + b + c + " " + s
	case HistoricalEventAddHfSiteLinkLinkType_PrisonAbstractBuilding:
		return h + " was imprisoned in " + b + c + " " + s
	case HistoricalEventAddHfSiteLinkLinkType_PrisonSiteBuildingProfile:
		return h + " was imprisoned in " + b + c + " " + s
	case HistoricalEventAddHfSiteLinkLinkType_SeatOfPower:
		return h + " ruled from " + b + c + " " + s
	default:
		return h + " LINKED TO " + s
	}
}

func (x *HistoricalEventAgreementFormed) Html() string { // TODO
	return "UNKNWON HistoricalEventAgreementFormed"
}

func (x *HistoricalEventAgreementMade) Html() string { // TODO
	return "UNKNWON HistoricalEventAgreementMade"
}

func (x *HistoricalEventAgreementRejected) Html() string { // TODO
	return "UNKNWON HistoricalEventAgreementRejected"
}

func (x *HistoricalEventArtifactClaimFormed) Html() string {
	a := artifact(x.ArtifactId)
	switch x.Claim {
	case HistoricalEventArtifactClaimFormedClaim_Heirloom:
		return a + " was made a family heirloom by " + hf(x.HistFigureId)
	case HistoricalEventArtifactClaimFormedClaim_Symbol:
		p := world.Entities[x.EntityId].Position(x.PositionProfileId).Name_
		e := entity(x.EntityId)
		return a + " was made a symbol of the " + p + " by " + e
	case HistoricalEventArtifactClaimFormedClaim_Treasure:
		c := ""
		if x.Circumstance != HistoricalEventArtifactClaimFormedCircumstance_Unknown {
			c = " " + x.Circumstance.String()
		}
		if x.HistFigureId != -1 {
			return a + " was claimed by " + hf(x.HistFigureId) + c
		} else if x.EntityId != -1 {
			return a + " was claimed by " + entity(x.EntityId) + c
		}
	}
	return a + " was claimed"
}

func (x *HistoricalEventArtifactCopied) Html() string {
	s := util.If(x.FromOriginal, "made a copy of the original", "aquired a copy of")
	return fmt.Sprintf("%s %s %s %s of %s, keeping it%s",
		entity(x.DestEntityId), s, artifact(x.ArtifactId), siteStructure(x.SourceSiteId, x.SourceStructureId, "from"),
		entity(x.SourceEntityId), siteStructure(x.DestSiteId, x.DestStructureId, "within"))
}

func (x *HistoricalEventArtifactCreated) Html() string {
	a := artifact(x.ArtifactId)
	h := hf(x.HistFigureId)
	s := ""
	if x.SiteId != -1 {
		s = site(x.SiteId, " in ")
	}
	if !x.NameOnly {
		return h + " created " + a + s
	}
	c := ""
	if x.Circumstance != nil {
		switch x.Circumstance.Type {
		case HistoricalEventArtifactCreatedCircumstanceType_Defeated:
			c = " after defeating " + hf(x.Circumstance.Defeated)
		case HistoricalEventArtifactCreatedCircumstanceType_Favoritepossession:
			c = " as the item was a favorite possession"
		case HistoricalEventArtifactCreatedCircumstanceType_Preservebody:
			c = " by preserving part of the body"
		}
	}
	switch x.Reason {
	case HistoricalEventArtifactCreatedReason_SanctifyHf:
		return fmt.Sprintf("%s received its name%s from %s in order to sanctify %s%s", a, s, h, hf(x.SanctifyHf), c)
	default:
		return fmt.Sprintf("%s received its name%s from %s %s", a, s, h, c)
	}
}

func (x *HistoricalEventArtifactDestroyed) Html() string {
	return fmt.Sprintf("%s was destroyed by %s in %s", artifact(x.ArtifactId), entity(x.DestroyerEnid), site(x.SiteId, ""))
}

func (x *HistoricalEventArtifactFound) Html() string {
	w := ""
	if x.SiteId != -1 {
		w = site(x.SiteId, "")
		if x.SitePropertyId != -1 {
			w = property(x.SiteId, x.SitePropertyId) + " in " + w
		}
	}
	return fmt.Sprintf("%s was found in %s by %s", artifact(x.ArtifactId), w, util.If(x.HistFigureId != -1, hf(x.HistFigureId), "an unknown creature"))
}

func (x *HistoricalEventArtifactGiven) Html() string {
	r := ""
	if x.ReceiverHistFigureId != -1 {
		r = hf(x.ReceiverHistFigureId)
		if x.ReceiverEntityId != -1 {
			r += " of " + entity(x.ReceiverEntityId)
		}
	} else if x.ReceiverEntityId != -1 {
		r += entity(x.ReceiverEntityId)
	}
	g := ""
	if x.GiverHistFigureId != -1 {
		g = hf(x.GiverHistFigureId)
		if x.GiverEntityId != -1 {
			g += " of " + entity(x.GiverEntityId)
		}
	} else if x.GiverEntityId != -1 {
		g += entity(x.GiverEntityId)
	}
	reason := ""
	switch x.Reason {
	case HistoricalEventArtifactGivenReason_PartOfTradeNegotiation:
		reason = " as part of a trade negotiation"
	}
	return fmt.Sprintf("%s was offered to %s by %s%s", artifact(x.ArtifactId), r, g, reason)
}
func (x *HistoricalEventArtifactLost) Html() string {
	w := ""
	if x.SubregionId != -1 {
		w = region(x.SubregionId)
	}
	if x.SiteId != -1 {
		w = site(x.SiteId, "")
		if x.SitePropertyId != -1 {
			w = property(x.SiteId, x.SitePropertyId) + " in " + w
		}
	}
	return fmt.Sprintf("%s was lost in %s", artifact(x.ArtifactId), w)
}

func (x *HistoricalEventArtifactPossessed) Html() string {
	a := artifact(x.ArtifactId)
	h := hf(x.HistFigureId)
	w := ""
	if x.SubregionId != -1 {
		w = region(x.SubregionId)
	}
	if x.SiteId != -1 {
		w = site(x.SiteId, "")
	}
	c := ""
	switch x.Circumstance {
	case HistoricalEventArtifactPossessedCircumstance_HfIsDead:
		c = " after the death of " + hf(x.CircumstanceId)
	}

	switch x.Reason {
	case HistoricalEventArtifactPossessedReason_ArtifactIsHeirloomOfFamilyHfid:
		return fmt.Sprintf("%s was aquired in %s by %s as an heirloom of %s%s", a, w, h, hf(x.ReasonId), c)
	case HistoricalEventArtifactPossessedReason_ArtifactIsSymbolOfEntityPosition:
		return fmt.Sprintf("%s was aquired in %s by %s as a symbol of authority within %s%s", a, w, h, entity(x.ReasonId), c)
	}
	return fmt.Sprintf("%s was claimed in %s by %s%s", a, w, h, c) // TODO wording
}

func (x *HistoricalEventArtifactRecovered) Html() string {
	a := artifact(x.ArtifactId)
	h := hf(x.HistFigureId)
	w := ""
	if x.SubregionId != -1 {
		w = "in " + region(x.SubregionId)
	}
	if x.SiteId != -1 {
		w = site(x.SiteId, "in ")
		if x.StructureId != -1 {
			w = siteStructure(x.SiteId, x.StructureId, "from")
		}
	}
	return fmt.Sprintf("%s was recovered %s by %s", a, w, h)
}

func (x *HistoricalEventArtifactStored) Html() string { // TODO export siteProperty
	if x.HistFigureId != -1 {
		return fmt.Sprintf("%s stored %s in %s", hf(x.HistFigureId), artifact(x.ArtifactId), site(x.SiteId, ""))
	} else {
		return fmt.Sprintf("%s was stored in %s", artifact(x.ArtifactId), site(x.SiteId, ""))
	}
}

func (x *HistoricalEventArtifactTransformed) Html() string {
	return fmt.Sprintf("%s was made from %s by %s in %s", artifact(x.NewArtifactId), artifact(x.OldArtifactId), hf(x.HistFigureId), site(x.SiteId, "")) // TODO wording
}

func (x *HistoricalEventAssumeIdentity) Html() string {
	h := hf(x.TricksterHfid)
	i := identity(x.IdentityId)
	if x.TargetEnid == -1 {
		return fmt.Sprintf(`%s assumed the identity "%s"`, h, i)
	} else {
		return fmt.Sprintf(`%s fooled %s into believing %s was "%s"`, h, entity(x.TargetEnid), pronoun(x.TricksterHfid), i)
	}
}

func (x *HistoricalEventAttackedSite) Html() string {
	atk := entity(x.AttackerCivId)
	def := siteCiv(x.SiteCivId, x.DefenderCivId)
	generals := ""
	if x.AttackerGeneralHfid != -1 {
		generals += ". " + util.Capitalize(hf(x.AttackerGeneralHfid)) + " led the attack"
		if x.DefenderGeneralHfid != -1 {
			generals += ", and the defenders were led by " + hf(x.DefenderGeneralHfid)
		}
	}
	mercs := ""
	if x.AttackerMercEnid != -1 {
		mercs += fmt.Sprintf(". %s were hired by the attackers", entity(x.AttackerMercEnid))
	}
	if x.ASupportMercEnid != -1 {
		mercs += fmt.Sprintf(". %s were hired as scouts by the attackers", entity(x.ASupportMercEnid))
	}
	if x.DefenderMercEnid != -1 {
		mercs += fmt.Sprintf(". The defenders hired %s", entity(x.DefenderMercEnid))
	}
	if x.DSupportMercEnid != -1 {
		mercs += fmt.Sprintf(". The defenders hired %s as scouts", entity(x.DSupportMercEnid))
	}
	return fmt.Sprintf("%s attacked %s at %s%s%s", atk, def, site(x.SiteId, ""), generals, mercs)
}

func (x *HistoricalEventBodyAbused) Html() string {
	s := "the " + util.If(len(x.Bodies) > 1, "bodies", "body") + " of " + hfList(x.Bodies) + " " + util.If(len(x.Bodies) > 1, "were", "was")

	switch x.AbuseType {
	case HistoricalEventBodyAbusedAbuseType_Animated:
		s += " animated" + util.If(x.Histfig != -1, " by "+hf(x.Histfig), "") + site(x.SiteId, " in ")
	case HistoricalEventBodyAbusedAbuseType_Flayed:
		s += " flayed and the skin stretched over " + structure(x.SiteId, x.Structure) + " by " + entity(x.Civ) + site(x.SiteId, " in ")
	case HistoricalEventBodyAbusedAbuseType_Hung:
		s += " hung from a tree by " + entity(x.Civ) + site(x.SiteId, " in ")
	case HistoricalEventBodyAbusedAbuseType_Impaled:
		s += " impaled on " + articled(x.ItemMat+" "+x.ItemSubtype.String()) + " by " + entity(x.Civ) + site(x.SiteId, " in ")
	case HistoricalEventBodyAbusedAbuseType_Mutilated:
		s += " horribly mutilated by " + entity(x.Civ) + site(x.SiteId, " in ")
	case HistoricalEventBodyAbusedAbuseType_Piled:
		s += " added to a "
		switch x.PileType {
		case HistoricalEventBodyAbusedPileType_Grislymound:
			s += "grisly mound"
		case HistoricalEventBodyAbusedPileType_Grotesquepillar:
			s += "grotesque pillar"
		case HistoricalEventBodyAbusedPileType_Gruesomesculpture:
			s += "gruesome sculpture"
		}
		s += " by " + entity(x.Civ) + site(x.SiteId, " in ")
	}

	return s
}

func (x *HistoricalEventBuildingProfileAcquired) Html() string {
	return util.If(x.AcquirerEnid != -1, entity(x.AcquirerEnid), hf(x.AcquirerHfid)) +
		util.If(x.PurchasedUnowned, " purchased ", " inherited ") +
		property(x.SiteId, x.BuildingProfileId) + site(x.SiteId, " in") +
		util.If(x.LastOwnerHfid != -1, " formerly owned by "+hf(x.LastOwnerHfid), "")
}

func (x *HistoricalEventCeremony) Html() string {
	r := entity(x.CivId) + " held a ceremony in " + site(x.SiteId, "")
	if e, ok := world.Entities[x.CivId]; ok {
		o := e.Occasion[x.OccasionId]
		r += " as part of " + o.Name()
		s := o.Schedule[x.ScheduleId]
		if len(s.Feature) > 0 {
			r += ". The event featured " + andList(util.Map(s.Feature, feature))
		}
	}
	return r
}

func (x *HistoricalEventChangeHfBodyState) Html() string {
	r := hf(x.Hfid)
	switch x.BodyState {
	case HistoricalEventChangeHfBodyStateBodyState_EntombedAtSite:
		r += " was entombed"
	}
	if x.StructureId != -1 {
		r += " within " + structure(x.SiteId, x.StructureId)
	}
	r += site(x.SiteId, " in ")
	return r
}

func (x *HistoricalEventChangeHfJob) Html() string {
	w := ""
	if x.SubregionId != -1 {
		w = " in " + region(x.SubregionId)
	}
	if x.SiteId != -1 {
		w = " in " + site(x.SiteId, "")
	}
	old := articled(strcase.ToDelimited(x.OldJob, ' '))
	new := articled(strcase.ToDelimited(x.NewJob, ' '))
	if x.OldJob == "standard" {
		return hf(x.Hfid) + " became " + new + w
	} else if x.NewJob == "standard" {
		return hf(x.Hfid) + " stopped being " + old + w
	} else {
		return hf(x.Hfid) + " gave up being " + old + " to become a " + new + w
	}
}

func (x *HistoricalEventChangeHfState) Html() string {
	r := ""
	switch x.Reason {
	case HistoricalEventChangeHfStateReason_BeWithMaster:
		r = " in order to be with the master"
	case HistoricalEventChangeHfStateReason_ConvictionExile, HistoricalEventChangeHfStateReason_ExiledAfterConviction:
		r = " after being exiled following a criminal conviction"
	case HistoricalEventChangeHfStateReason_FailedMood:
		r = " after failing to create an artifact"
	case HistoricalEventChangeHfStateReason_Flight:
	case HistoricalEventChangeHfStateReason_GatherInformation:
		r = " to gather information"
	case HistoricalEventChangeHfStateReason_GreatDealOfStress:
		r = " after a great deal of stress" // TODO check
	case HistoricalEventChangeHfStateReason_LackOfSleep:
		r = " after a lack of sleep" // TODO check
	case HistoricalEventChangeHfStateReason_OnAPilgrimage:
		r = " on a pilgrimage"
	case HistoricalEventChangeHfStateReason_Scholarship:
		r = " in order to pursue scholarship"
	case HistoricalEventChangeHfStateReason_UnableToLeaveLocation:
		r = " after being unable to leave the location" // TODO check
	}

	switch x.State {
	case HistoricalEventChangeHfStateState_Refugee:
		return hf(x.Hfid) + " fled " + location(x.SiteId, "to", x.SubregionId, "into")
	case HistoricalEventChangeHfStateState_Settled:
		switch x.Reason {
		case HistoricalEventChangeHfStateReason_BeWithMaster, HistoricalEventChangeHfStateReason_Scholarship:
			return hf(x.Hfid) + " moved to study " + site(x.SiteId, "in") + r
		case HistoricalEventChangeHfStateReason_Flight:
			return hf(x.Hfid) + " fled " + site(x.SiteId, "to")
		case HistoricalEventChangeHfStateReason_ConvictionExile, HistoricalEventChangeHfStateReason_ExiledAfterConviction:
			return hf(x.Hfid) + " departed " + site(x.SiteId, "to") + r
		case HistoricalEventChangeHfStateReason_None:
			return hf(x.Hfid) + " settled " + location(x.SiteId, "in", x.SubregionId, "in")
		}
	case HistoricalEventChangeHfStateState_Visiting:
		return hf(x.Hfid) + " visited " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateState_Wandering:
		if x.SubregionId != -1 {
			return hf(x.Hfid) + " began wandering " + region(x.SubregionId)
		} else {
			return hf(x.Hfid) + " began wandering the wilds"
		}
	}

	switch x.Mood {
	case HistoricalEventChangeHfStateMood_Berserk:
		return hf(x.Hfid) + " went berserk " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Catatonic:
		return hf(x.Hfid) + " stopped responding to the outside world " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Fell:
		return hf(x.Hfid) + " was taken by a fell mood " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Fey:
		return hf(x.Hfid) + " was taken by a fey mood " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Insane:
		return hf(x.Hfid) + " became crazed " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Macabre:
		return hf(x.Hfid) + " began to skulk and brood " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Melancholy:
		return hf(x.Hfid) + " was striken by melancholy " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Possessed:
		return hf(x.Hfid) + " was posessed " + site(x.SiteId, "in") + r
	case HistoricalEventChangeHfStateMood_Secretive:
		return hf(x.Hfid) + " withdrew from society " + site(x.SiteId, "in") + r
	}
	return "UNKNWON HistoricalEventChangeHfState"
}

func (x *HistoricalEventChangedCreatureType) Html() string {
	return hf(x.ChangerHfid) + " changed " + hf(x.ChangeeHfid) + " from " + articled(x.OldRace) + " to " + articled(x.NewRace)
}

func (x *HistoricalEventCompetition) Html() string {
	e := world.Entities[x.CivId]
	o := e.Occasion[x.OccasionId]
	s := o.Schedule[x.ScheduleId]
	return entity(x.CivId) + " held a " + strcase.ToDelimited(s.Type.String(), ' ') + site(x.SiteId, " in") + " as part of the " + o.Name() +
		". Competing " + util.If(len(x.CompetitorHfid) > 1, "were ", "was ") + hfList(x.CompetitorHfid) + ". " +
		util.Capitalize(hf(x.WinnerHfid)) + " was the victor"
}

func (x *HistoricalEventCreateEntityPosition) Html() string {
	c := entity(x.Civ)
	if x.SiteCiv != x.Civ {
		c = entity(x.SiteCiv) + " of " + c
	}
	if x.Histfig != -1 {
		c = hf(x.Histfig) + " of " + c
	} else {
		c = "members of " + c
	}
	switch x.Reason {
	case HistoricalEventCreateEntityPositionReason_AsAMatterOfCourse:
		return c + " created the position of " + x.Position + " as a matter of course"
	case HistoricalEventCreateEntityPositionReason_Collaboration:
		return c + " collaborated to create the position of " + x.Position
	case HistoricalEventCreateEntityPositionReason_ForceOfArgument:
		return c + " created the position of " + x.Position + " trough force of argument"
	case HistoricalEventCreateEntityPositionReason_ThreatOfViolence:
		return c + " compelled the creation of the position of " + x.Position + " with threats of violence"
	case HistoricalEventCreateEntityPositionReason_WaveOfPopularSupport:
		return c + " created the position of " + x.Position + ", pushed by a wave of popular support"
	}
	return c + " created the position of " + x.Position
}

func (x *HistoricalEventCreatedSite) Html() string {
	f := util.If(x.ResidentCivId != -1, " for "+entity(x.ResidentCivId), "")
	if x.BuilderHfid != -1 {
		return hf(x.BuilderHfid) + " created " + site(x.SiteId, "") + f
	}
	return siteCiv(x.SiteCivId, x.CivId) + " founded " + site(x.SiteId, "") + f

}

func (x *HistoricalEventCreatedStructure) Html() string { // TODO rebuild/rebuilt
	if x.BuilderHfid != -1 {
		return hf(x.BuilderHfid) + " thrust a spire of slade up from the underworld, naming it " + structure(x.SiteId, x.StructureId) +
			", and established a gateway between worlds in " + site(x.SiteId, "")
	}
	return siteCiv(x.SiteCivId, x.CivId) + util.If(x.Rebuilt, " rebuild ", " constructed ") + siteStructure(x.SiteId, x.StructureId, "")
}

func (x *HistoricalEventCreatedWorldConstruction) Html() string {
	return siteCiv(x.SiteCivId, x.CivId) + " finished the contruction of " + worldConstruction(x.Wcid) +
		" connecting " + site(x.SiteId1, "") + " with " + site(x.SiteId2, "") +
		util.If(x.MasterWcid != -1, " as part of "+worldConstruction(x.MasterWcid), "")
}

func (x *HistoricalEventCreatureDevoured) Html() string {
	return hf(x.Eater) + " devoured " + util.If(x.Victim != -1, hf(x.Victim), articled(x.Race)) +
		util.If(x.Entity != -1, " of "+entity(x.Entity), "") +
		location(x.SiteId, " in", x.SubregionId, " in")
}

func (x *HistoricalEventDanceFormCreated) Html() string {
	reason := ""
	switch x.Reason {
	case HistoricalEventDanceFormCreatedReason_GlorifyHf:
		reason = " in order to glorify " + hf(x.ReasonId)
	}
	circumstance := ""
	switch x.Circumstance {
	case HistoricalEventDanceFormCreatedCircumstance_Dream:
		circumstance = " after a dream"
	case HistoricalEventDanceFormCreatedCircumstance_DreamAboutHf:
		circumstance = " after a dreaming about " + hf(x.CircumstanceId)
	case HistoricalEventDanceFormCreatedCircumstance_Nightmare:
		circumstance = " after a nightmare"
	case HistoricalEventDanceFormCreatedCircumstance_PrayToHf:
		circumstance = " after praying to " + hf(x.CircumstanceId)
	}
	return danceForm(x.FormId) + " was created by " + hf(x.HistFigureId) + location(x.SiteId, " in", x.SubregionId, " in") + reason + circumstance
}

func (x *HistoricalEventDestroyedSite) Html() string { // TODO NoDefeatMention
	return entity(x.AttackerCivId) + " defeated " + siteCiv(x.SiteCivId, x.DefenderCivId) + " and destroyed " + site(x.SiteId, "")
}

func (x *HistoricalEventDiplomatLost) Html() string { // TODO
	return "UNKNWON HistoricalEventDiplomatLost"
}

func (x *HistoricalEventEntityAllianceFormed) Html() string {
	return entityList(x.JoiningEnid) + " swore to support " + entity(x.InitiatingEnid) + " in war if the latter did likewise"
}

func (x *HistoricalEventEntityBreachFeatureLayer) Html() string {
	return siteCiv(x.SiteEntityId, x.CivEntityId) + " breached the Underworld at " + site(x.SiteId, "")
}

func (x *HistoricalEventEntityCreated) Html() string {
	if x.CreatorHfid != -1 {
		return hf(x.CreatorHfid) + " formed " + entity(x.EntityId) + siteStructure(x.SiteId, x.StructureId, "in")
	} else {
		return entity(x.EntityId) + " formed" + siteStructure(x.SiteId, x.StructureId, "in")
	}
}

func (x *HistoricalEventEntityDissolved) Html() string {
	return entity(x.EntityId) + " dissolved after " + x.Reason.String()
}

func (x *HistoricalEventEntityEquipmentPurchase) Html() string { // todo check hfid
	return entity(x.EntityId) + " purchased " + equipmentLevel(x.NewEquipmentLevel) + " equipment"
}

func (x *HistoricalEventEntityExpelsHf) Html() string {
	return "UNKNWON HistoricalEventEntityExpelsHf"
}

func (x *HistoricalEventEntityFledSite) Html() string {
	return "UNKNWON HistoricalEventEntityFledSite"
}

func (x *HistoricalEventEntityIncorporated) Html() string { // TODO site
	return entity(x.JoinerEntityId) + util.If(x.PartialIncorporation, " began operating at the direction of ", " fully incorporated into ") +
		entity(x.JoinedEntityId) + " under the leadership of " + hf(x.LeaderHfid)
}

func (x *HistoricalEventEntityLaw) Html() string {
	switch x.LawAdd {
	case HistoricalEventEntityLawLawAdd_Harsh:
		return hf(x.HistFigureId) + " laid a series of oppressive edicts upon " + entity(x.EntityId)
	}
	switch x.LawRemove {
	case HistoricalEventEntityLawLawRemove_Harsh:
		return hf(x.HistFigureId) + " lifted numerous oppressive laws from " + entity(x.EntityId)
	}
	return hf(x.HistFigureId) + " UNKNOWN LAW upon " + entity(x.EntityId)
}

func (x *HistoricalEventEntityOverthrown) Html() string {
	return hf(x.InstigatorHfid) + " toppled the government of " + util.If(x.OverthrownHfid != -1, hf(x.OverthrownHfid)+" of ", "") + entity(x.EntityId) + " and " +
		util.If(x.PosTakerHfid == x.InstigatorHfid, "assumed control", "placed "+hf(x.PosTakerHfid)+" in power") + site(x.SiteId, " in") +
		util.If(len(x.ConspiratorHfid) > 0, ". The support of "+hfList(x.ConspiratorHfid)+" was crucial to the coup", "")
}

func (x *HistoricalEventEntityPersecuted) Html() string {
	var l []string
	if len(x.ExpelledHfid) > 0 {
		l = append(l, hfList(x.ExpelledHfid)+util.If(len(x.ExpelledHfid) > 1, " were", " was")+" expelled")
	}
	if len(x.PropertyConfiscatedFromHfid) > 0 {
		l = append(l, "most property was confiscated")
	}
	if x.DestroyedStructureId != -1 {
		l = append(l, structure(x.SiteId, x.DestroyedStructureId)+" was destroyed"+util.If(x.ShrineAmountDestroyed > 0, " along with several smaller sacred sites", ""))
	} else if x.ShrineAmountDestroyed > 0 {
		l = append(l, "some sacred sites were desecrated")
	}
	return hf(x.PersecutorHfid) + " of " + entity(x.PersecutorEnid) + " persecuted " + entity(x.TargetEnid) + " in " + site(x.SiteId, "") +
		util.If(len(l) > 0, ". "+util.Capitalize(andList(l)), "")
}

func (x *HistoricalEventEntityPrimaryCriminals) Html() string { // TODO structure
	switch x.Action {
	case HistoricalEventEntityPrimaryCriminalsAction_EntityPrimaryCriminals:
		return entity(x.EntityId) + " became the primary criminal organization in " + site(x.SiteId, "")
	}
	return "UNKNWON HistoricalEventEntityPrimaryCriminals"
}

func (x *HistoricalEventEntityRampagedInSite) Html() string { // TODO
	return "UNKNWON HistoricalEventEntityRampagedInSite"
}

func (x *HistoricalEventEntityRelocate) Html() string {
	switch x.Action {
	case HistoricalEventEntityRelocateAction_EntityRelocate:
		return entity(x.EntityId) + " moved" + siteStructure(x.SiteId, x.StructureId, "to")
	}
	return "UNKNWON HistoricalEventEntityRelocate"
}

func (x *HistoricalEventEntitySearchedSite) Html() string { // TODO
	return "UNKNWON HistoricalEventEntitySearchedSite"
}

func (x *HistoricalEventFailedFrameAttempt) Html() string {
	return hf(x.FramerHfid) + " attempted to frame " + hf(x.TargetHfid) + " for " + x.Crime.String() +
		util.If(x.PlotterHfid != -1, " at the behest of "+hf(x.PlotterHfid), "") +
		" by fooling " + hf(x.FooledHfid) + " and " + entity(x.ConvicterEnid) +
		" with fabricated evidence, but nothing came of it"
}

func (x *HistoricalEventFailedIntrigueCorruption) Html() string {
	action := ""
	switch x.Action {
	case HistoricalEventFailedIntrigueCorruptionAction_BribeOfficial:
		action = "have law enforcement look the other way"
	case HistoricalEventFailedIntrigueCorruptionAction_BringIntoNetwork:
		action = "have someone to act on plots and schemes"
	case HistoricalEventFailedIntrigueCorruptionAction_CorruptInPlace:
		action = "have an agent"
	case HistoricalEventFailedIntrigueCorruptionAction_InduceToEmbezzle:
		action = "secure embezzled funds"
	}
	method := ""
	switch x.Method {
	case HistoricalEventFailedIntrigueCorruptionMethod_BlackmailOverEmbezzlement:
		method = "made a blackmail threat, due to embezzlement using the position " + position(x.RelevantEntityId, x.RelevantPositionProfileId, x.CorruptorHfid) + " of " + entity(x.RelevantEntityId)
	case HistoricalEventFailedIntrigueCorruptionMethod_Bribe:
		method = "offered a bribe"
	case HistoricalEventFailedIntrigueCorruptionMethod_Flatter:
		method = "made flattering remarks"
	case HistoricalEventFailedIntrigueCorruptionMethod_Intimidate:
		method = "made a threat"
	case HistoricalEventFailedIntrigueCorruptionMethod_OfferImmortality:
		method = "offered immortality"
	case HistoricalEventFailedIntrigueCorruptionMethod_Precedence:
		method = "pulled rank as " + position(x.RelevantEntityId, x.RelevantPositionProfileId, x.CorruptorHfid) + " of " + entity(x.RelevantEntityId)
	case HistoricalEventFailedIntrigueCorruptionMethod_ReligiousSympathy:
		method = "played for sympathy" + util.If(x.RelevantIdForMethod != -1, " by appealing to shared worship of "+hf(x.RelevantIdForMethod), "")
	case HistoricalEventFailedIntrigueCorruptionMethod_RevengeOnGrudge:
		method = "offered revenge upon the persecutor " + hf(x.RelevantIdForMethod)
	}
	fail := "The plan failed"
	switch x.TopValue {
	case HistoricalEventFailedIntrigueCorruptionTopValue_Law:
		fail = hf(x.TargetHfid) + " valued the law and refused"
	case HistoricalEventFailedIntrigueCorruptionTopValue_Power:
	}
	switch x.TopFacet {
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Ambition:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_AnxietyPropensity:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Confidence:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_EnvyPropensity:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Fearlessness:
		fail += ", despite being afraid"
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Greed:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Hope:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Pride:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_StressVulnerability:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Swayable:
		fail += ", despite being swayed by the emotional appeal" // TODO
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Vanity:
	case HistoricalEventFailedIntrigueCorruptionTopFacet_Vengeful:
	}
	return hf(x.CorruptorHfid) + " attempted to corrupt " + hf(x.TargetHfid) +
		" in order to " + action + location(x.SiteId, " in", x.SubregionId, " in") + ". " +
		util.Capitalize(util.If(x.LureHfid != -1,
			hf(x.LureHfid)+" lured "+hfShort(x.TargetHfid)+" to a meeting with "+hfShort(x.CorruptorHfid)+", where the latter",
			hfShort(x.CorruptorHfid)+" met with "+hfShort(x.TargetHfid))) +
		util.If(x.FailedJudgmentTest, ", while completely misreading the situation,", "") + " " + method + ". " + fail
}

func (x *HistoricalEventFieldBattle) Html() string {
	atk := entity(x.AttackerCivId)
	def := entity(x.DefenderCivId)
	generals := ""
	if x.AttackerGeneralHfid != -1 {
		generals += ". " + util.Capitalize(hf(x.AttackerGeneralHfid)) + " led the attack"
		if x.DefenderGeneralHfid != -1 {
			generals += ", and the defenders were led by " + hf(x.DefenderGeneralHfid)
		}
	}
	mercs := ""
	if x.AttackerMercEnid != -1 {
		mercs += fmt.Sprintf(". %s were hired by the attackers", entity(x.AttackerMercEnid))
	}
	if x.ASupportMercEnid != -1 {
		mercs += fmt.Sprintf(". %s were hired as scouts by the attackers", entity(x.ASupportMercEnid))
	}
	if x.DefenderMercEnid != -1 {
		mercs += fmt.Sprintf(". The defenders hired %s", entity(x.DefenderMercEnid))
	}
	if x.DSupportMercEnid != -1 {
		mercs += fmt.Sprintf(". The defenders hired %s as scouts", entity(x.DSupportMercEnid))
	}
	return fmt.Sprintf("%s attacked %s at %s%s%s", atk, def, region(x.SubregionId), generals, mercs)
}

func (x *HistoricalEventFirstContact) Html() string { // TODO
	return "UNKNWON HistoricalEventFirstContact"
}

func (x *HistoricalEventGamble) Html() string {
	outcome := ""
	switch d := x.NewAccount - x.OldAccount; {
	case d <= -5000:
		outcome = "lost a fortune"
	case d <= -1000:
		outcome = "did poorly"
	case d <= 1000:
		outcome = "did well"
	case d <= 5000:
		outcome = "made a fortune"
	}
	return hf(x.GamblerHfid) + " " + outcome + " gambling" + siteStructure(x.SiteId, x.StructureId, " in") +
		util.If(x.OldAccount >= 0 && x.NewAccount < 0, " and went into debt", "")
}

func (x *HistoricalEventHfAbducted) Html() string {
	return hf(x.TargetHfid) + " was abducted " + location(x.SiteId, "from", x.SubregionId, "from") + " by " + hf(x.SnatcherHfid)
}

func (x *HistoricalEventHfAttackedSite) Html() string {
	return hf(x.AttackerHfid) + " attacked " + siteCiv(x.SiteCivId, x.DefenderCivId) + site(x.SiteId, " in")
}

func (x *HistoricalEventHfConfronted) Html() string {
	return hf(x.Hfid) + " aroused " + x.Situation.String() + location(x.SiteId, " in", x.SubregionId, " in") + " after " +
		andList(util.Map(x.Reason, func(r HistoricalEventHfConfrontedReason) string {
			switch r {
			case HistoricalEventHfConfrontedReason_Ageless:
				return " appearing not to age"
			case HistoricalEventHfConfrontedReason_Murder:
				return "a murder"
			}
			return ""
		}))
}
func (x *HistoricalEventHfConvicted) Html() string { // TODO no_prison_available, beating, hammerstrokes, interrogator_hfid
	r := util.If(x.ConfessedAfterApbArrestEnid != -1, "after being recognized and arrested, ", "")
	switch {
	case x.SurveiledCoconspirator:
		r += "due to ongoing surveillance on a coconspiratior, " + hf(x.CoconspiratorHfid) + ", as the plot unfolded, "
	case x.SurveiledContact:
		r += "due to ongoing surveillance on the contact " + hf(x.ContactHfid) + " as the plot unfolded, "
	case x.SurveiledConvicted:
		r += "due to ongoing surveillance as the plot unfolded, "
	case x.SurveiledTarget:
		r += "due to ongoing surveillance on the target " + hf(x.TargetHfid) + " as the plot unfolded, "
	}
	r += hf(x.ConvictedHfid) + util.If(x.ConfessedAfterApbArrestEnid != -1, " confessed and", "") + " was " + util.If(x.WrongfulConviction, "wrongfully ", "") + "convicted " +
		util.If(x.ConvictIsContact, "as a go-between in a conspiracy to commit ", "of ") + x.Crime.String() + " by " + entity(x.ConvicterEnid)
	if x.FooledHfid != -1 {
		r += " after " + hf(x.FramerHfid) + " fooled " + hf(x.FooledHfid) + " with fabricated evidence" +
			util.If(x.PlotterHfid != -1, " at the behest of "+hf(x.PlotterHfid), "")
	}
	if x.CorruptConvicterHfid != -1 {
		r += " and the corrupt " + hf(x.CorruptConvicterHfid) + " through the machinations of " + hf(x.PlotterHfid)
	}
	switch {
	case x.DeathPenalty:
		r += " and sentenced to death"
	case x.Exiled:
		r += " and exiled"
	case x.PrisonMonths > 0:
		r += fmt.Sprintf(" and imprisoned for a term of %d years", x.PrisonMonths/12)
	}
	if x.HeldFirmInInterrogation {
		r += ". " + hfShort(x.ConvictedHfid) + " revealed nothing during interrogation"
	} else if len(x.ImplicatedHfid) > 0 {
		r += ". " + hfShort(x.ConvictedHfid) + " implicated " + hfList(x.ImplicatedHfid) + " during interrogation" +
			util.If(x.DidNotRevealAllInInterrogation, " but did not reveal eaverything", "")
	}
	return r
}

func (x *HistoricalEventHfDestroyedSite) Html() string {
	return hf(x.AttackerHfid) + " routed " + siteCiv(x.SiteCivId, x.DefenderCivId) + " and destroyed " + site(x.SiteId, "")
}

func (x *HistoricalEventHfDied) Html() string { // TODO force cause enum
	return "UNKNWON HistoricalEventHfDied"
}

func (x *HistoricalEventHfDisturbedStructure) Html() string {
	return hf(x.HistFigId) + " disturbed " + siteStructure(x.SiteId, x.StructureId, "")
}

func (x *HistoricalEventHfDoesInteraction) Html() string {
	i := strings.Index(x.InteractionAction, " ")
	if i > 0 {
		return hf(x.DoerHfid) + " " + x.InteractionAction[:i+1] + hf(x.TargetHfid) + x.InteractionAction[i:] + util.If(x.Site != -1, site(x.Site, " in"), "")
	} else {
		return hf(x.DoerHfid) + " UNKNOWN INTERACTION " + hf(x.TargetHfid) + util.If(x.Site != -1, site(x.Site, " in"), "")
	}
}

func (x *HistoricalEventHfEnslaved) Html() string {
	return hf(x.SellerHfid) + " sold " + hf(x.EnslavedHfid) + " to " + entity(x.PayerEntityId) + site(x.MovedToSiteId, " in")
}

func (x *HistoricalEventHfEquipmentPurchase) Html() string { // TODO site, structure, region
	return hf(x.GroupHfid) + " purchased " + equipmentLevel(x.Quality) + " equipment"
}

func (x *HistoricalEventHfFreed) Html() string { // TODO
	return "UNKNWON HistoricalEventHfFreed"
}

func (x *HistoricalEventHfGainsSecretGoal) Html() string {
	switch x.SecretGoal {
	case HistoricalEventHfGainsSecretGoalSecretGoal_Immortality:
		return hf(x.Hfid) + " became obsessed with " + posessivePronoun(x.Hfid) + " own mortality and sought to extend " + posessivePronoun(x.Hfid) + " life by any means"
	}
	return hf(x.Hfid) + " UNKNOWN SECRET GOAL"
}

func (x *HistoricalEventHfInterrogated) Html() string { // TODO
	return "UNKNWON HistoricalEventHfInterrogated"
}

func (x *HistoricalEventHfLearnsSecret) Html() string { return "UNKNWON HistoricalEventHfLearnsSecret" }
func (x *HistoricalEventHfNewPet) Html() string       { return "UNKNWON HistoricalEventHfNewPet" }
func (x *HistoricalEventHfPerformedHorribleExperiments) Html() string {
	return "UNKNWON HistoricalEventHfPerformedHorribleExperiments"
}
func (x *HistoricalEventHfPrayedInsideStructure) Html() string {
	return "UNKNWON HistoricalEventHfPrayedInsideStructure"
}
func (x *HistoricalEventHfPreach) Html() string { return "UNKNWON HistoricalEventHfPreach" }
func (x *HistoricalEventHfProfanedStructure) Html() string {
	return "UNKNWON HistoricalEventHfProfanedStructure"
}
func (x *HistoricalEventHfRansomed) Html() string    { return "UNKNWON HistoricalEventHfRansomed" }
func (x *HistoricalEventHfReachSummit) Html() string { return "UNKNWON HistoricalEventHfReachSummit" }
func (x *HistoricalEventHfRecruitedUnitTypeForEntity) Html() string {
	return "UNKNWON HistoricalEventHfRecruitedUnitTypeForEntity"
}
func (x *HistoricalEventHfRelationshipDenied) Html() string {
	return "UNKNWON HistoricalEventHfRelationshipDenied"
}
func (x *HistoricalEventHfReunion) Html() string { return "UNKNWON HistoricalEventHfReunion" }
func (x *HistoricalEventHfRevived) Html() string { return "UNKNWON HistoricalEventHfRevived" }
func (x *HistoricalEventHfSimpleBattleEvent) Html() string {
	return "UNKNWON HistoricalEventHfSimpleBattleEvent"
}
func (x *HistoricalEventHfTravel) Html() string { return "UNKNWON HistoricalEventHfTravel" }
func (x *HistoricalEventHfViewedArtifact) Html() string {
	return "UNKNWON HistoricalEventHfViewedArtifact"
}
func (x *HistoricalEventHfWounded) Html() string { return "UNKNWON HistoricalEventHfWounded" }
func (x *HistoricalEventHfsFormedIntrigueRelationship) Html() string {
	return "UNKNWON HistoricalEventHfsFormedIntrigueRelationship"
}
func (x *HistoricalEventHfsFormedReputationRelationship) Html() string {
	return "UNKNWON HistoricalEventHfsFormedReputationRelationship"
}
func (x *HistoricalEventHolyCityDeclaration) Html() string {
	return "UNKNWON HistoricalEventHolyCityDeclaration"
}
func (x *HistoricalEventInsurrectionStarted) Html() string {
	return "UNKNWON HistoricalEventInsurrectionStarted"
}
func (x *HistoricalEventItemStolen) Html() string { return "UNKNWON HistoricalEventItemStolen" }
func (x *HistoricalEventKnowledgeDiscovered) Html() string {
	return "UNKNWON HistoricalEventKnowledgeDiscovered"
}
func (x *HistoricalEventMasterpieceArchConstructed) Html() string {
	return "UNKNWON HistoricalEventMasterpieceArchConstructed"
}
func (x *HistoricalEventMasterpieceEngraving) Html() string {
	return "UNKNWON HistoricalEventMasterpieceEngraving"
}
func (x *HistoricalEventMasterpieceFood) Html() string {
	return "UNKNWON HistoricalEventMasterpieceFood"
}
func (x *HistoricalEventMasterpieceItem) Html() string {
	return "UNKNWON HistoricalEventMasterpieceItem"
}
func (x *HistoricalEventMasterpieceItemImprovement) Html() string {
	return "UNKNWON HistoricalEventMasterpieceItemImprovement"
}
func (x *HistoricalEventMasterpieceLost) Html() string {
	return "UNKNWON HistoricalEventMasterpieceLost"
}
func (x *HistoricalEventMerchant) Html() string { return "UNKNWON HistoricalEventMerchant" }
func (x *HistoricalEventModifiedBuilding) Html() string {
	return "UNKNWON HistoricalEventModifiedBuilding"
}
func (x *HistoricalEventMusicalFormCreated) Html() string {
	return "UNKNWON HistoricalEventMusicalFormCreated"
}
func (x *HistoricalEventNewSiteLeader) Html() string { return "UNKNWON HistoricalEventNewSiteLeader" }
func (x *HistoricalEventPeaceAccepted) Html() string { return "UNKNWON HistoricalEventPeaceAccepted" }
func (x *HistoricalEventPeaceRejected) Html() string { return "UNKNWON HistoricalEventPeaceRejected" }
func (x *HistoricalEventPerformance) Html() string   { return "UNKNWON HistoricalEventPerformance" }
func (x *HistoricalEventPlunderedSite) Html() string { return "UNKNWON HistoricalEventPlunderedSite" }
func (x *HistoricalEventPoeticFormCreated) Html() string {
	return "UNKNWON HistoricalEventPoeticFormCreated"
}
func (x *HistoricalEventProcession) Html() string     { return "UNKNWON HistoricalEventProcession" }
func (x *HistoricalEventRazedStructure) Html() string { return "UNKNWON HistoricalEventRazedStructure" }
func (x *HistoricalEventReclaimSite) Html() string    { return "UNKNWON HistoricalEventReclaimSite" }
func (x *HistoricalEventRegionpopIncorporatedIntoEntity) Html() string {
	return "UNKNWON HistoricalEventRegionpopIncorporatedIntoEntity"
}
func (x *HistoricalEventRemoveHfEntityLink) Html() string {
	return "UNKNWON HistoricalEventRemoveHfEntityLink"
}
func (x *HistoricalEventRemoveHfHfLink) Html() string { return "UNKNWON HistoricalEventRemoveHfHfLink" }
func (x *HistoricalEventRemoveHfSiteLink) Html() string {
	return "UNKNWON HistoricalEventRemoveHfSiteLink"
}
func (x *HistoricalEventReplacedStructure) Html() string {
	return "UNKNWON HistoricalEventReplacedStructure"
}
func (x *HistoricalEventSiteDied) Html() string    { return "UNKNWON HistoricalEventSiteDied" }
func (x *HistoricalEventSiteDispute) Html() string { return "UNKNWON HistoricalEventSiteDispute" }
func (x *HistoricalEventSiteRetired) Html() string { return "UNKNWON HistoricalEventSiteRetired" }
func (x *HistoricalEventSiteSurrendered) Html() string {
	return "UNKNWON HistoricalEventSiteSurrendered"
}
func (x *HistoricalEventSiteTakenOver) Html() string { return "UNKNWON HistoricalEventSiteTakenOver" }
func (x *HistoricalEventSiteTributeForced) Html() string {
	return "UNKNWON HistoricalEventSiteTributeForced"
}
func (x *HistoricalEventSneakIntoSite) Html() string { return "UNKNWON HistoricalEventSneakIntoSite" }
func (x *HistoricalEventSpottedLeavingSite) Html() string {
	return "UNKNWON HistoricalEventSpottedLeavingSite"
}
func (x *HistoricalEventSquadVsSquad) Html() string { return "UNKNWON HistoricalEventSquadVsSquad" }
func (x *HistoricalEventTacticalSituation) Html() string {
	return "UNKNWON HistoricalEventTacticalSituation"
}
func (x *HistoricalEventTrade) Html() string { return "UNKNWON HistoricalEventTrade" }
func (x *HistoricalEventWrittenContentComposed) Html() string {
	return "UNKNWON HistoricalEventWrittenContentComposed"
}

func (x *HistoricalEventAgreementConcluded) Html() string {
	return "UNKNWON HistoricalEventAgreementConcluded"
}
func (x *HistoricalEventMasterpieceDye) Html() string {
	return "UNKNWON HistoricalEventMasterpieceDye"
}

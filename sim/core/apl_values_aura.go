package core

import (
	"fmt"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type APLValueAuraIsActive struct {
	defaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraIsActive(config *proto.APLValueAuraIsActive) APLValue {
	aura := rot.aplGetAura(rot.getSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraIsActive{
		aura: aura,
	}
}
func (value *APLValueAuraIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraIsActive) GetBool(sim *Simulation) bool {
	return value.aura.Get().IsActive()
}
func (value *APLValueAuraIsActive) String() string {
	return fmt.Sprintf("Aura Is Active(%s)", value.aura.String())
}

type APLValueAuraRemainingTime struct {
	defaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraRemainingTime(config *proto.APLValueAuraRemainingTime) APLValue {
	aura := rot.aplGetAura(rot.getSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraRemainingTime{
		aura: aura,
	}
}
func (value *APLValueAuraRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAuraRemainingTime) GetDuration(sim *Simulation) time.Duration {
	return value.aura.Get().RemainingDuration(sim)
}
func (value *APLValueAuraRemainingTime) String() string {
	return fmt.Sprintf("Aura Remaining Time(%s)", value.aura.String())
}

type APLValueAuraNumStacks struct {
	defaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraNumStacks(config *proto.APLValueAuraNumStacks) APLValue {
	aura := rot.aplGetAura(rot.getSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	if aura.Get().MaxStacks == 0 {
		rot.validationWarning("%s is not a stackable aura", ProtoToActionID(config.AuraId))
		return nil
	}
	return &APLValueAuraNumStacks{
		aura: aura,
	}
}
func (value *APLValueAuraNumStacks) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueAuraNumStacks) GetInt(sim *Simulation) int32 {
	return value.aura.Get().GetStacks()
}
func (value *APLValueAuraNumStacks) String() string {
	return fmt.Sprintf("Aura Num Stacks(%s)", value.aura.String())
}

type APLValueAuraInternalCooldown struct {
	defaultAPLValueImpl
	aura AuraReference
}

func (rot *APLRotation) newValueAuraInternalCooldown(config *proto.APLValueAuraInternalCooldown) APLValue {
	aura := rot.aplGetICDAura(rot.getSourceUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLValueAuraInternalCooldown{
		aura: aura,
	}
}
func (value *APLValueAuraInternalCooldown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueAuraInternalCooldown) GetDuration(sim *Simulation) time.Duration {
	return value.aura.Get().Icd.TimeToReady(sim)
}
func (value *APLValueAuraInternalCooldown) String() string {
	return fmt.Sprintf("Aura ICD(%s)", value.aura.String())
}

type APLValueAuraShouldRefresh struct {
	defaultAPLValueImpl
	aura       AuraReference
	maxOverlap APLValue
}

func (rot *APLRotation) newValueAuraShouldRefresh(config *proto.APLValueAuraShouldRefresh) APLValue {
	aura := rot.aplGetAura(rot.getTargetUnit(config.SourceUnit), config.AuraId)
	if aura.Get() == nil {
		return nil
	}

	maxOverlap := rot.coerceTo(rot.newAPLValue(config.MaxOverlap), proto.APLValueType_ValueTypeDuration)
	if maxOverlap == nil {
		maxOverlap = rot.newValueConst(&proto.APLValueConst{Val: "0ms"})
	}

	return &APLValueAuraShouldRefresh{
		aura:       aura,
		maxOverlap: maxOverlap,
	}
}
func (value *APLValueAuraShouldRefresh) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAuraShouldRefresh) GetBool(sim *Simulation) bool {
	return value.aura.Get().ShouldRefreshExclusiveEffects(sim, value.maxOverlap.GetDuration(sim))
}
func (value *APLValueAuraShouldRefresh) String() string {
	return fmt.Sprintf("Should Refresh Aura(%s)", value.aura.String())
}

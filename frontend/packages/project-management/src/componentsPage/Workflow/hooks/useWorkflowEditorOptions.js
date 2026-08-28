/**
 * 工作流编辑器：枚举 → Select / 颜色文案（含 i18n label）。
 */
'use client';

import { useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { localizeSelectOptions } from '@repo/project-management/enum/task';
import {
	WORKFLOW_STATE_COLOR_HEX,
	WORKFLOW_STATE_COLOR_LABEL_KEYS,
	WORKFLOW_STATE_COLOR_OPTIONS,
	WORKFLOW_TRANSITION_TRIGGER_OPTIONS,
	getWorkflowPhaseI18nKey,
} from '@repo/project-management/enum/workflow';
import { useOptionsStore } from '@repo/project-management/store/options';

export function useWorkflowEditorOptions() {
	const t = useTranslations();
	const stageOptionsRaw = useOptionsStore((state) => state.stageOptions);

	const stageOptions = useMemo(
		() => localizeSelectOptions(stageOptionsRaw, t, getWorkflowPhaseI18nKey),
		[stageOptionsRaw, t],
	);

	const triggerOptions = useMemo(
		() =>
			WORKFLOW_TRANSITION_TRIGGER_OPTIONS.map((opt) => ({
				value: opt.value,
				label: t(opt.labelKey),
			})),
		[t],
	);

	const colorOptions = useMemo(
		() =>
			WORKFLOW_STATE_COLOR_OPTIONS.map((value) => ({
				value,
				label: t(
					WORKFLOW_STATE_COLOR_LABEL_KEYS[value] || WORKFLOW_STATE_COLOR_LABEL_KEYS.gray,
				),
				hex: WORKFLOW_STATE_COLOR_HEX[value],
			})),
		[t],
	);

	return { stageOptions, triggerOptions, colorOptions };
}

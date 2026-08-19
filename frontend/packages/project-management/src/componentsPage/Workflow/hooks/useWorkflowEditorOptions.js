/**
 * 工作流编辑器：枚举 → Select / 色块选项（含 i18n label）。
 */
'use client';

import { useMemo } from 'react';
import { useTranslations } from 'next-intl';
import {
	WORKFLOW_STATE_COLOR_HEX,
	WORKFLOW_STATE_COLOR_LABEL_KEYS,
	WORKFLOW_STATE_COLOR_OPTIONS,
	WORKFLOW_STATE_TYPE_OPTIONS,
	WORKFLOW_TRANSITION_MEMBER_OPTIONS,
	WORKFLOW_TRANSITION_TRIGGER_OPTIONS,
} from '@repo/project-management/enum/workflow';

export function useWorkflowEditorOptions() {
	const t = useTranslations();

	const colorMap = useMemo(
		() => Object.fromEntries(Object.entries(WORKFLOW_STATE_COLOR_HEX)),
		[],
	);

	const typeOptions = useMemo(
		() =>
			WORKFLOW_STATE_TYPE_OPTIONS.map((opt) => ({
				value: opt.value,
				label: t(opt.labelKey),
			})),
		[t],
	);

	const triggerOptions = useMemo(
		() =>
			WORKFLOW_TRANSITION_TRIGGER_OPTIONS.map((opt) => ({
				value: opt.value,
				label: t(opt.labelKey),
			})),
		[t],
	);

	const memberOptions = useMemo(
		() =>
			WORKFLOW_TRANSITION_MEMBER_OPTIONS.map((opt) => ({
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

	return { colorMap, typeOptions, triggerOptions, memberOptions, colorOptions };
}

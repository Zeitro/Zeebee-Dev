'use client';

import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Button, Progress, message } from 'antd';
import { CloudSyncOutlined } from '@ant-design/icons';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import CategoryContentPanel from './CategoryContentPanel';
import { buildDocumentCategoriesFromChecklist, getChecklistProgress } from './checklist';
import DocumentCategoryButton from './DocumentCategoryButton';
// 同步是异步任务：POST 立即返回 job，之后轮询 status 接口直到终态。
// 代理（undici）约 30s 会掐断长请求，所以绝不能让 HTTP 请求一直等同步结果。
const SYNC_POLL_INTERVAL = 2000;
const SYNC_POLL_MAX_MS = 10 * 60 * 1000; // 与后端 driveSyncJobTimeout 对齐
const SYNC_TERMINAL_STATUSES = ['Success', 'PartialFailed', 'Failed'];
const PARSE_POLL_INTERVAL = 5000; // 每 5 秒请求一次 parse/result
const PARSE_POLL_MAX_MS = 60 * 60 * 1000; // 最长轮询 1 小时，超时则停止并提示
/** 解析任务终态；Running 为进行中，Finished / Failed 为终态 */
const PARSE_TERMINAL_STATUSES = ['Finished', 'Failed'];
const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const DocumentsTab = ({ rowData, onRefresh }) => {
	// createRequest() 每次渲染返回新的 request 引用，useMemo 固定一次，避免回调失稳
	const { request } = useMemo(() => createRequest(), []);
	const [syncing, setSyncing] = useState(false);
	const [messageApi, contextHolder] = message.useMessage();
	// 组件卸载后停止轮询 / 状态更新，避免切换菜单回来时遗留的异步写入已卸载组件
	const mountedRef = useRef(true);
	useEffect(() => {
		mountedRef.current = true;
		return () => {
			mountedRef.current = false;
		};
	}, []);
	// onRefresh 以 ref 持有，供同步完成后刷新文档列表（拉取最新 drive_synced_at），
	// 同时避免把易变的 onRefresh 直接放进轮询/恢复 effect 的依赖里造成重复触发。
	const onRefreshRef = useRef(onRefresh);
	useEffect(() => {
		onRefreshRef.current = onRefresh;
	}, [onRefresh]);

	/**
	 * 轮询 GET /borrower_property_transaction/{id}/parse/result?request_id= 直到终态。
	 * 必须用 getState() 读最新 params：updateParseFileParams 是同步写入，
	 * 但 useCallback 闭包里的 parseFileParams 仍是旧快照，会导致刚写入却读不到。
	 * @param {number|string} transactionId
	 * @param {number|string} documentId
	 * @returns {Promise<object|null>}
	 */
	const pollParseStatus = useCallback(
		async (transactionId, documentId, params) => {
			if (!transactionId || !documentId) return null;

			const deadline = Date.now() + PARSE_POLL_MAX_MS;
			let result = null;

			while (!result || !PARSE_TERMINAL_STATUSES.includes(result?.status)) {
				if (!mountedRef.current || Date.now() > deadline) {
					return null;
				}

				const query = new URLSearchParams({
					...params,
					type: 'request_id',
				}).toString();
				result = await request(
					`/borrower_property_transaction/${transactionId}/parse/result?${query}`,
					{ method: 'GET' },
				);

				if (PARSE_TERMINAL_STATUSES.includes(result?.status)) {
					if (result?.status === 'Finished') {
						messageApi.success('AI parsing completed');
					} else {
						messageApi.error('AI parsing failed');
					}
					return result;
				}

				await sleep(PARSE_POLL_INTERVAL);
			}

			return result;
		},
		[messageApi, request],
	);

	// 轮询任务状态直到终态；shouldStop 为真或组件卸载时提前退出。返回终态结果或 null。
	const pollSyncStatus = useCallback(
		async (transactionId, initial) => {
			const deadline = Date.now() + SYNC_POLL_MAX_MS;
			let result = initial;
			while (!result || !SYNC_TERMINAL_STATUSES.includes(result?.status)) {
				if (!mountedRef.current) {
					return null;
				}
				if (Date.now() > deadline) {
					message.warning(
						'Sync is taking longer than expected. Please check again later.',
					);
					return null;
				}
				await sleep(SYNC_POLL_INTERVAL);
				if (!mountedRef.current) {
					return null;
				}
				result = await request(
					`/borrower_property_transaction/${transactionId}/sync_to_drive/status`,
					{ method: 'GET' },
				);
			}
			return result;
		},
		[request],
	);

	const runSync = useCallback(async () => {
		const transactionId = rowData?.id;
		// syncing 时直接返回：防止在状态尚未刷新的瞬间重复触发
		if (!transactionId || syncing) {
			return;
		}
		try {
			setSyncing(true);
			// 先检查 Google Drive 是否已连接
			const status = await request('/drive/auth_status', { method: 'GET' });
			if (!status?.is_authorized) {
				message.warning(
					'Google Drive is not connected. Please connect it in Setting > System first.',
				);
				return;
			}
			// 触发异步任务，立即返回 { status: 'Running', job_id }。
			// 后端对正在运行的任务幂等：若已有 Running 任务会返回同一个，不会重复同步。
			const started = await request(
				`/borrower_property_transaction/${transactionId}/sync_to_drive`,
				{ method: 'POST' },
			);
			message.info('Sync started. Syncing to Google Drive in the background…');

			const result = await pollSyncStatus(transactionId, started);
			// On terminal status, refresh the document list so newly synced files show
			// the Google Drive tag immediately. The full result (counts / folder link /
			// per-file list) is recorded in the transaction Log instead of a modal.
			// Failed means the whole run failed with nothing persisted, so skip refresh.
			if (result && result.status !== 'Failed') {
				onRefreshRef.current?.();
			}
		} catch (error) {
			message.error('Sync failed: ' + (error?.message || 'unknown error'));
		} finally {
			if (mountedRef.current) {
				setSyncing(false);
			}
		}
	}, [request, rowData?.id, syncing, pollSyncStatus]);

	// On returning to this page, if a Running job still exists, restore the loading
	// state and keep polling so the button isn't shown as clickable (avoiding a
	// duplicate sync); refresh the document list on completion to show the tags.
	useEffect(() => {
		const transactionId = rowData?.id;
		if (!transactionId) {
			return;
		}
		let active = true;
		(async () => {
			try {
				const current = await request(
					`/borrower_property_transaction/${transactionId}/sync_to_drive/status`,
					{ method: 'GET' },
				);
				if (!active || !mountedRef.current) {
					return;
				}
				if (current?.status === 'Running') {
					setSyncing(true);
					const result = await pollSyncStatus(transactionId, current);
					if (active && mountedRef.current && result && result.status !== 'Failed') {
						onRefreshRef.current?.();
					}
				}
			} catch {
				// 状态探测为尽力而为，失败时静默忽略
			} finally {
				if (active && mountedRef.current) {
					setSyncing(false);
				}
			}
		})();
		return () => {
			active = false;
		};
	}, [rowData?.id, request, pollSyncStatus]);

	const categories = useMemo(
		() => buildDocumentCategoriesFromChecklist(rowData?.checklist),
		[rowData?.checklist],
	);
	const { completed, total } = useMemo(
		() => getChecklistProgress(rowData?.checklist),
		[rowData?.checklist],
	);
	const progressPercent = total > 0 ? Math.round((completed / total) * 100) : 0;

	const [activeKey, setActiveKey] = useState(null);

	// rowData 异步加载后 categories 才有值；useState 初值不会随 default 更新，需派生当前选中项
	const resolvedActiveKey = useMemo(() => {
		if (!categories.length) return null;
		if (activeKey && categories.some((c) => c.key === activeKey)) {
			return activeKey;
		}
		const firstPending = categories.find((c) => c.pendingCount > 0 && !c.completed);
		return firstPending?.key ?? categories[0]?.key;
	}, [activeKey, categories]);

	const activeCategory = categories.find((c) => c.key === resolvedActiveKey) ?? categories[0];
	const requiredCategories = categories.filter((c) => c.is_required);
	const optionalCategories = categories.filter((c) => !c.is_required);

	return (
		<div className="flex gap-4 p-4">
			{contextHolder}
			<div className="flex w-[311px] shrink-0 flex-col gap-2">
				<div className="flex gap-2">
					{/* Import 按钮本期不需要，先隐藏 */}
					<Button
						className="flex-1"
						icon={<CloudSyncOutlined />}
						loading={syncing}
						onClick={runSync}
					>
						Sync to Google Drive
					</Button>
				</div>

				<div className="rounded-lg border border-base-border bg-background p-6">
					<div className="mb-2 flex items-center gap-8">
						<p className="m-0 text-base font-medium leading-6 text-foreground">
							Overall progress
						</p>
						<p className="m-0 text-sm leading-[22px] text-base-gray">
							{total > 0 ? `${completed}/${total}` : displayValue(null)}
						</p>
					</div>
					<Progress
						percent={progressPercent}
						showInfo={false}
						strokeColor="var(--blue-dark)"
						trailColor="var(--blue-light-1)"
						size="small"
					/>
				</div>

				<div className="rounded-lg border border-base-border bg-background py-4">
					<p className="m-0 px-4 pb-2 text-sm font-medium leading-[22px] text-base-gray">
						Required
					</p>
					{requiredCategories.map((cat) => (
						<DocumentCategoryButton
							key={cat.key}
							category={cat}
							isActive={cat.key === resolvedActiveKey}
							onSelect={setActiveKey}
						/>
					))}

					<div className="my-2 border-y border-base-border px-4 py-2">
						<p className="m-0 text-sm font-medium leading-[22px] text-base-gray">
							Optional
						</p>
					</div>

					{optionalCategories.map((cat) => (
						<DocumentCategoryButton
							key={cat.key}
							category={cat}
							isActive={cat.key === resolvedActiveKey}
							onSelect={setActiveKey}
							isHideStatusTag={true}
						/>
					))}
				</div>
			</div>

			<CategoryContentPanel
				category={activeCategory}
				onRefresh={onRefresh}
				transactionId={rowData?.id}
				pollParseStatus={pollParseStatus}
			/>
		</div>
	);
};

export default DocumentsTab;

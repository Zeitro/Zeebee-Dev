'use client';

import React, { useMemo } from 'react';
import CollapsibleDocumentList from '@/components/CollapsibleDocumentList';
import { NoData } from '@repo/shared-methods';
import { toDocumentListItems } from '../constants';
import { transactionDocumentColumns } from './documentColumns';
import styles from './index.module.scss';

function SectionLabel({ title, count }) {
	return (
		<div className={styles.sectionLabel}>
			<p className={styles.sectionLabelTitle}>{title}</p>
			<span className={styles.sectionLabelCount}>· {count}</span>
		</div>
	);
}

const UploadDocumentList = ({ category }) => {
	const columns = useMemo(() => transactionDocumentColumns(), []);

	const { requiredList, optionalList } = useMemo(() => {
		const documents = category?.documents ?? [];
		return {
			requiredList: toDocumentListItems(documents.filter((d) => d.required)),
			optionalList: toDocumentListItems(documents.filter((d) => !d.required)),
		};
	}, [category?.documents]);

	if (!requiredList.length && !optionalList.length) {
		return (
			<div className="px-6 py-6">
				<NoData />
			</div>
		);
	}

	const hasRequired = requiredList.length > 0;
	const hasOptional = optionalList.length > 0;

	return (
		<div className="px-6 py-5">
			{hasRequired && (
				<section className="mt-4">
					<SectionLabel title="Required" count={requiredList.length} />
					<CollapsibleDocumentList columns={columns} dataSource={requiredList} />
				</section>
			)}

			{hasOptional && (
				<section className={hasRequired ? 'mt-8' : 'mt-4'}>
					<SectionLabel title="Optional" count={optionalList.length} />
					<CollapsibleDocumentList columns={columns} dataSource={optionalList} />
				</section>
			)}
		</div>
	);
};

export default UploadDocumentList;

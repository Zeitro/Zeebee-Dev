'use client';
import React, { useState } from 'react';
import LoanFilterTable from '.';
import { pendingColumns } from '../column';
import { transformArrayStructure } from '@/lib/common';
import { loanPendingStatus } from '@/enum/loan';
const PendingPage = () => {
    const [type, setType] = useState('');
    const [searchValue, setSearchValue] = useState('');
    const [loanPurpose, setLoanPurpose] = useState([]);
    const [loanType, setLoanType] = useState([]);
    const [lender,setLender] = useState([]);
    const [status,setStatus] = useState([]);
	return (
		<LoanFilterTable
			columns={pendingColumns}
			scroll={{ x: 3500, y: 'calc(100% - 54px)' }}
			dataSource={[{ status: 'quoting', id: 1 }]}
			filterList={[
				{
					type: 'select-input',
                    selectValue: type,
                    value: searchValue,
					options: [
						{
							label: 'Address',
							value: 'address',
						},
						{
							label: 'State',
							value: 'state',
						},
						{
							label: 'City',
							value: 'city',
						},
						{
							label: 'County',
							value: 'county',
						},
						{
							label: 'Zip Code',
							value: 'zipCode',
						},
					],
				},
				{
					type: 'treeSelect',
					label: 'Loan purpose',
                    popupMatchSelectWidth: 172,
                    value:loanPurpose,
					treeData: [
						{
							title: 'Purchase',
							value: 'purchase',
							key: 'purchase',
							isLeaf: true,
						},
						{
							title: 'Rate-term refi',
							value: 'rateTermRefi',
							key: 'rateTermRefi',
							isLeaf: true,
						},
						{
							title: 'Cash-out refi',
							value: 'cashOutRefi',
							key: 'cashOutRefi',
							isLeaf: true,
						},
					],
				},
				{
					type: 'treeSelect',
					label: 'Loan type',
                    popupMatchSelectWidth: 195,
                    value: loanType,
					treeData: [
						{
							title: 'Bridge',
							value: 'bridge',
							key: 'bridge',
							isLeaf: true,
						},
						{
							title: 'Rehab',
							value: 'rehab',
							key: 'rehab',
							isLeaf: true,
						},
						{
							title: 'Ground up construction',
							value: 'groundUpConstruction',
							key: 'groundUpConstruction',
							isLeaf: true,
						},
					],
				},
				{
					type: 'treeSelect',
					label: 'Lender',
                    popupMatchSelectWidth:172,
                    value: lender,
					treeData: [
						{
							title: 'Bridge',
							value: 'bridge',
							key: 'bridge',
							isLeaf: true,
						},
						{
							title: 'Rehab',
							value: 'rehab',
							key: 'rehab',
							isLeaf: true,
						},
						{
							title: 'Ground up construction',
							value: 'groundUpConstruction',
							key: 'groundUpConstruction',
							isLeaf: true,
						},
					],
				},
				{
					type: 'treeSelect',
					label: 'Status',
                    value:status,
                    popupMatchSelectWidth:172,
					treeData: transformArrayStructure({
						sourceArray: loanPendingStatus,
						outputKey: 'title',
					}).map((i) => ({ ...i, key: i.value })),
				},
			]}
		/>
	);
};
export default PendingPage;

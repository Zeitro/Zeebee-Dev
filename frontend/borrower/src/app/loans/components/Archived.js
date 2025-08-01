'use client';
import React, { useEffect, useMemo, useState } from 'react';
import FilterTable from '@/components/FilterTable';
import { archivedColumns } from './column';
import { transformArrayStructure } from '@/lib/common';
import {useLoanStore} from '@/store/loan'
const ArchivedLoans = ({hasFilterField,hasPagination,filterWrapClassName}) => {
    const {  loanTypeOptions } = useLoanStore();
    const [type, setType] = useState('');
    const [searchValue, setSearchValue] = useState('');
    const [loanPurpose, setLoanPurpose] = useState([]);
    const [loanType, setLoanType] = useState([]);
    const [lender, setLender] = useState([]);
    const [status, setStatus] = useState([]);
    const filterList = useMemo(
        () => [
            {
                type: 'select-input',
                selectValue: type,
                value: searchValue,
                selectProps: {
                    onChange: (e) => setType(e),
                },
                onChange: ({ target: { value } }) => setSearchValue(value),
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
				label: 'Loan type',
				popupMatchSelectWidth: 195,
				value: loanType,
				placeholder: 'Loan type',
				onChange: (value) => setLoanType(value),
				treeData: transformArrayStructure({
					sourceArray: loanTypeOptions,
					outputKey: 'title',
				}).map((i) => ({ ...i, key: i.value })),
			}, 
            {
                type: 'treeSelect',
                label: 'Lender',
                placeholder: 'Lender',
                popupMatchSelectWidth: 172,
                value: lender,
                onChange: (value) => setLender(value),
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
            
        ],
        [
            searchValue,
            type,
            loanPurpose,
            loanType,
            lender,
            status,
        ],
    );

    return (
        <FilterTable
            columns={archivedColumns}
            scroll={{ x: 2000, y: 'calc(100% - 54px)' }}
            dataSource={[{ status: 'quoting', id: 1 }]}
            filterList={filterList}
            hasFilterField={hasFilterField}
            pagination={hasPagination}
            filterWrapClassName={filterWrapClassName}
        />
    );
};
export default ArchivedLoans;

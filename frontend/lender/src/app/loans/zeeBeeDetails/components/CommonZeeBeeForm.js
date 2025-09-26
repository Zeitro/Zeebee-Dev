'use client';
import React from 'react';
import ApplicantInfo from '../../components/CreateZeeBee/ApplicantInfo';
import BorrowerInfo from '../../components/CreateZeeBee/BorrowerInfo';
import PropertyInfo from '../../components/CreateZeeBee/PropertyInfo';
import Loan from '../../components/CreateZeeBee/Loan';
import Experience from '../../components/CreateZeeBee/Experience';
import ThirdPartyService from '../../components/CreateZeeBee/ThirdPartyService';
import Declaration from '../../components/CreateZeeBee/Declaration';
import Demographics from '../../components/CreateZeeBee/Demographics';
import ReviewAndSign from '../../components/CreateZeeBee/ReviewAndSign';
import Quoting from './Quoting';

const CommonZeeBeeForm = ({ form, pageKey, children, setIsChange, rowData }) => {
	const getContent = () => {
		switch (pageKey) {
			case 'borrowerInformation':
				return (
					<div className="px-6 py-6">
						<ApplicantInfo form={form} />
						<BorrowerInfo form={form} />
						{children}
					</div>
				);
			case 'propertyAndLoan':
				return (
					<div className="px-6 py-6">
						<PropertyInfo form={form} />
						<Loan form={form} />
						{children}
					</div>
				);
			case 'experienceAndThirdService':
				return (
					<div className="px-6 py-6">
						<Experience form={form} rowData={rowData} setIsChange={setIsChange} />
						<ThirdPartyService
							form={form}
							rowData={rowData}
							setIsChange={setIsChange}
						/>
						{children}
					</div>
				);
			case 'declarationAndDemographics':
				return (
					<div className="px-6 py-6">
						<Declaration form={form} />
						<Demographics form={form} />
						{children}
					</div>
				);
			case 'reviewAndSign':
				return (
					<div className="px-6 py-6">
						<ReviewAndSign form={form} />
						{children}
					</div>
				);
			case 'Quoting':
				return (
					<div>
						<Quoting />
						{children}
					</div>
				);
			default:
				return '';
		}
	};
	return getContent();
};
export default CommonZeeBeeForm;

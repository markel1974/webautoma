package executor

var stub = `
{
	"id": "a346ff6a-b967-4ab9-a50b-97da5a495b2a",
	"version": "2.0",
	"timeout": 10000,
	"name": "h2d",
	"quit": true,
	"id": "03571100282",
	"url": "https://tuttodati.it/",
	"tests": [
		{
			"commands": [
				{
					"id": "test111",
					"comment": "",
					"command": "enableDebug",
					"target": "",
					"targets": [],
					"value": "0"
				},
				{
					"id": "f5cfe0a7-f929-4a50-84e7-6745ad3874ba",
					"comment": "",
					"command": "click",
					"target": "name=search",
					"targets": [
						[
							"name=search",
							"name"
						],
						[
							"css=.border-2",
							"css:finder"
						],
						[
							"xpath=//input[@name='search']",
							"xpath:attributes"
						],
						[
							"xpath=//input",
							"xpath:position"
						]
					],
					"value": ""
				},
				{
					"id": "test2",
					"comment": "",
					"command": "humanWait",
					"target": "",
					"targets": [],
					"value": "10"
				},
				{
					"id": "69a27bf0-f96c-4231-b5c8-cfec2abe1602",
					"comment": "",
					"command": "type",
					"target": "name=search",
					"targets": [
						[
							"name=search",
							"name"
						],
						[
							"css=.border-2",
							"css:finder"
						],
						[
							"xpath=//input[@name='search']",
							"xpath:attributes"
						],
						[
							"xpath=//input",
							"xpath:position"
						]
					],
					"value": "{{id}}"
				},
				{
					"id": "test2",
					"comment": "",
					"command": "humanWait",
					"target": "",
					"targets": [],
					"value": "10"
				},
				{
					"id": "b0f230a3-5d64-4c85-bed4-fbf4c3959a1b",
					"comment": "",
					"command": "click",
					"target": "css=.bg-blue-500",
					"targets": [
						[
							"css=.bg-blue-500",
							"css:finder"
						],
						[
							"xpath=//button[@type='submit']",
							"xpath:attributes"
						],
						[
							"xpath=//button",
							"xpath:position"
						],
						[
							"xpath=//button[contains(.,'Cerca')]",
							"xpath:innerText"
						]
					],
					"value": ""
				},
				{
					"id": "test2",
					"comment": "",
					"command": "humanWait",
					"target": "",
					"targets": [],
					"value": "10"
				},
				{
					"id": "ede1c1c2-f699-42aa-ae23-2d28323e7679",
					"comment": "",
					"command": "click",
					"target": "xpath=//td/a",
					"targets": [
						[
							"linkText=ACQUA & SAPONE SRL",
							"linkText"
						],
						[
							"css=.p-2 > a",
							"css:finder"
						],
						[
							"xpath=//a[contains(text(),'ACQUA & SAPONE SRL')]",
							"xpath:link"
						],
						[
							"xpath=//a[contains(@href, '//tuttodati.it/aziende/57373-acqua-sapone-s-r-l')]",
							"xpath:href"
						],
						[
							"xpath=//td/a",
							"xpath:position"
						],
						[
							"xpath=//a[contains(.,'ACQUA & SAPONE SRL')]",
							"xpath:innerText"
						]
					],
					"value": ""
				},
				{
					"id": "12345",
					"comment": "",
					"command": "save",
					"target": "",
					"targets": [],
					"value": "{{id}}"
				}
			]
		}
	],
	"suites": [
		{
			"id": "1612f65c-347f-4f83-a55c-69f9b15e501d",
			"name": "Default Suite",
			"persistSession": false,
			"parallel": false,
			"timeout": 300,
			"tests": [
				"9d37e411-d3a1-432d-8e63-de105d3fd645"
			]
		}
	],
	"urls": [
		"https://sales.telecomitalia.local/"
	],
	"plugins": []
}
`

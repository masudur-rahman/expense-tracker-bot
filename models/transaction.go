package models

type TransactionType string

const (
	ExpenseTransaction  TransactionType = "Expense"
	IncomeTransaction   TransactionType = "Income"
	TransferTransaction TransactionType = "Transfer"
)

type Transaction struct {
	ID            int64
	Amount        float64
	SubcategoryID string
	Type          TransactionType
	SrcID         string
	DstID         string
	User          string
	Timestamp     int64
	Remarks       string
}

type TxnCategory struct {
	CatID string
	Name  string
}

type TxnSubcategory struct {
	SubCatID string
	Name     string
	CatID    string
}

var categories = []TxnCategory{
	{CatID: "food", Name: "Food"},
	{CatID: "house", Name: "Housing"},
	{CatID: "ent", Name: "Entertainment"},
	{CatID: "pc", Name: "Personal Care"},
	{CatID: "trv", Name: "Travel"},
	{CatID: "fin", Name: "Financial"},
	{CatID: "misc", Name: "Miscellaneous"},
}

var subcategories []TxnSubcategory

var foodSubs = []TxnSubcategory{
	{SubCatID: "food-rest", Name: "Restaurants", CatID: "food"},
	{SubCatID: "food-groc", Name: "Groceries", CatID: "food"},
	{SubCatID: "food-take", Name: "Takeout", CatID: "food"},
	{SubCatID: "food-snack", Name: "Snacks", CatID: "food"},
	{SubCatID: "food-fruit", Name: "Fruits", CatID: "food"},
}

var houseSubs = []TxnSubcategory{
	{SubCatID: "house-rent", Name: "Rent", CatID: "house"},
	{SubCatID: "house-util", Name: "Utilities", CatID: "house"},
	{SubCatID: "house-furn", Name: "Furniture", CatID: "house"},
	{SubCatID: "house-elec", Name: "Electronics", CatID: "house"},
	{SubCatID: "house-real", Name: "Real State", CatID: "house"},
}

var entSubs = []TxnSubcategory{
	{SubCatID: "ent-movie", Name: "Movies", CatID: "ent"},
	{SubCatID: "ent-sub", Name: "Subscription", CatID: "ent"},
	{SubCatID: "ent-rec", Name: "Recreation", CatID: "ent"},
	{SubCatID: "ent-books", Name: "Books", CatID: "ent"},
}

var pcSubs = []TxnSubcategory{
	{SubCatID: "pc-salon", Name: "Salon", CatID: "pc"},
	{SubCatID: "pc-toilet", Name: "Toiletries", CatID: "pc"},
	{SubCatID: "pc-gym", Name: "Gym", CatID: "pc"},
	{SubCatID: "pc-cloth", Name: "Clothing", CatID: "pc"},
	{SubCatID: "pc-health", Name: "Health", CatID: "pc"},
}

var trvSubs = []TxnSubcategory{
	{SubCatID: "trv-accom", Name: "Accommodation", CatID: "trv"},
	{SubCatID: "trv-dine", Name: "Dining", CatID: "trv"},
	{SubCatID: "trv-sight", Name: "Sightseeing", CatID: "trv"},
	{SubCatID: "trv-trans", Name: "Transportation", CatID: "trv"},
	{SubCatID: "trv-gift", Name: "Gifts", CatID: "trv"},
}

var finSubs = []TxnSubcategory{
	{SubCatID: "fin-sal", Name: "Salary", CatID: "fin"},
	{SubCatID: "fin-deposit", Name: "Deposit", CatID: "fin"},
	{SubCatID: "fin-with", Name: "Withdraw", CatID: "fin"},
	{SubCatID: "fin-dps", Name: "DPS", CatID: "fin"},
	{SubCatID: "fin-ccpay", Name: "Credit Card Payment", CatID: "fin"},
	{SubCatID: "fin-bank", Name: "Bank Transfer", CatID: "fin"},
	{SubCatID: "fin-loan", Name: "Loan", CatID: "fin"},
	{SubCatID: "fin-borrow", Name: "Borrow", CatID: "fin"},
	{SubCatID: "fin-tax", Name: "Tax", CatID: "fin"},
}

var miscSubs = []TxnSubcategory{
	{SubCatID: "misc-give", Name: "Giveaway", CatID: "misc"},
	{SubCatID: "misc-init", Name: "Initial Amount", CatID: "misc"},
	{SubCatID: "misc-misc", Name: "Misc", CatID: "misc"},
}

func init() {
	subcategories = append(subcategories, foodSubs...)
	subcategories = append(subcategories, houseSubs...)
	subcategories = append(subcategories, entSubs...)
	subcategories = append(subcategories, pcSubs...)
	subcategories = append(subcategories, trvSubs...)
	subcategories = append(subcategories, finSubs...)
	subcategories = append(subcategories, miscSubs...)
}

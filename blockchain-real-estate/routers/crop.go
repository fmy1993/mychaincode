/*
 * @Descripttion:
 * @version:
 * @Author: fmy1993
 * @Date: 2021-03-30 11:02:29
 * @LastEditors: fmy1993
 * @LastEditTime: 2021-05-12 21:59:33
 */
package routers

import (
	"encoding/json"
	"fmt"

	"github.com/fmy1993/BCexplorer/chaincode/blockchain-real-estate/lib"
	"github.com/fmy1993/BCexplorer/chaincode/blockchain-real-estate/utils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//增加数据				参数自动 [][]byte --> []string   pb.Response is a struct
func CreateCrop(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 3 { //2 { 增加参数
		return shim.Error("参数个数不满足")
	}
	datatype := args[0]
	id := args[1]
	hashinfo := args[2]

	if id == "" || hashinfo == "" || datatype == "" {
		return shim.Error("参数存在空值")
	}

	// 参数数据格式转换
	// var formattedTotalArea float64
	// if val, err := strconv.ParseFloat(totalArea, 64); err != nil {
	// 	return shim.Error(fmt.Sprintf("totalArea参数格式转换出错: %s", err))
	// } else {
	// 	formattedTotalArea = val
	// }

	//判断数据是否存在
	// resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{proprietor})
	// if err != nil || len(resultsProprietor) != 1 {
	// 	return shim.Error(fmt.Sprintf("业主proprietor信息验证失败%s", err))
	// }
	Crop := &lib.Crop{
		DataType: datatype,
		Id:       id,
		HashInfo: hashinfo,
	}
	// 写入账本   // []string{Crop.Id, Crop.HashInfo} 以两个字段做主键存入按一个字段查询查不到 ，复合主键的类型也必须对上
	if err := utils.WriteLedger(Crop, stub, lib.CropKey, []string{Crop.DataType + "-" + Crop.Id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	CropByte, err := json.Marshal(Crop)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(CropByte) //返回的数据会存在这个结构体的payload中
}

//查询上链数据列表
func QueryCrop(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var cropList []lib.Crop // crop数组
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.CropKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var crop lib.Crop // 声明了结构体
			err := json.Unmarshal(v, &crop)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryCropList-反序列化出错: %s", err))
			}
			cropList = append(cropList, crop)
		}
	}
	cropListByte, err := json.Marshal(cropList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryCropList-序列化出错: %s", err))
	}
	return shim.Success(cropListByte)
}

// 根据账本数据库的id更新账本数据库
func UpdateCrop(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 3 { //2 { 增加参数
		return shim.Error("参数个数不满足")
	}
	datatype := args[0]
	id := args[1]
	hashinfo := args[2]

	if id == "" || hashinfo == "" || datatype == "" {
		return shim.Error("参数存在空值")
	}
	Crop := &lib.Crop{
		DataType: datatype,
		Id:       id,
		HashInfo: hashinfo,
	}

	if err := utils.DelLedger(stub, lib.CropKey, []string{Crop.DataType + "-" + Crop.Id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	CropByte, err := json.Marshal(Crop)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(CropByte) //返回的数据会存在这个结构体的payload中
}

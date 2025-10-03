module.exports = (sequelize, DataTypes) => {
  const User = sequelize.define('User', {
    id: { type: DataTypes.UUID, defaultValue: DataTypes.UUIDV4, primaryKey: true },
    name: { type: DataTypes.STRING, allowNull: false },
    email: { type: DataTypes.STRING, unique: true, allowNull: false },
    password: { type: DataTypes.STRING, allowNull: false },
    role: { 
      type: DataTypes.ENUM('admin', 'staff', 'member'), 
      defaultValue: 'member', 
      allowNull: false 
    },
    // Khusus Member
    phoneNumber: { type: DataTypes.STRING, allowNull: true },
    address: { type: DataTypes.STRING, allowNull: true },
    packageId: { type: DataTypes.INTEGER, allowNull: true }, // Foreign Key ke GymPackage
    isActive: { type: DataTypes.BOOLEAN, defaultValue: true }, 
    refreshToken: { type: DataTypes.STRING(512), allowNull: true }, // Untuk JWT Refresh
  }, {});

  User.associate = function(models) {
    User.belongsTo(models.GymPackage, { foreignKey: 'packageId', as: 'gymPackage' });
    User.hasMany(models.Attendance, { foreignKey: 'userId', as: 'attendances' });
  };

  return User;
};